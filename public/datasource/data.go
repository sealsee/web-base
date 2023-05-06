package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sealsee/web-base/public/datasource/tx"
	"github.com/sealsee/web-base/public/setting"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var data *Data

type Data struct {
	masterDb    *sqlx.DB
	slaveDb     []*sqlx.DB
	redisDb     *redis.Client
	rabbitMQChn *amqp.Channel
	kfkProducer sarama.SyncProducer
	kfkConsumer sarama.ConsumerGroup
	mgoDb       *mongo.Database
}

func GetDB() *sqlx.DB {
	return data.masterDb
}

func GetSlaveDb() *sqlx.DB {
	return data.slaveDb[0]
}

func GetRedisClient() *redis.Client {
	return data.redisDb
}

func GetRabbitMQChn() *amqp.Channel {
	return data.rabbitMQChn
}

func GetKafkaPrd() sarama.SyncProducer {
	return data.kfkProducer
}

func GetKafkaCom() sarama.ConsumerGroup {
	return data.kfkConsumer
}

func GetMgoDb() *mongo.Database {
	return data.mgoDb
}

func InitCompent(d *setting.Datasource) (*Data, func(), error) {
	masterDb := newMasterDB(d.Master)
	slaveDb := newSlaveDB(d.Slave)
	redisClient := newRedis(d.Redis)
	rabbitMQChn := newRabbitMq(d.RabbitMQ)
	kfkProducer, kfkConsumer := newKafka(d.Kafka)
	mgoDb := newMongodb(d.Mongodb)

	cleanup := func() {
		masterDb.Close()
		for _, db := range slaveDb {
			db.Close()
		}
		redisClient.Close()
		if rabbitMQChn != nil {
			rabbitMQChn.Close()
		}
		if kfkProducer != nil {
			kfkProducer.Close()
		}
		if kfkConsumer != nil {
			kfkConsumer.Close()
		}
		if mgoDb != nil {
			mgoDb.Client().Disconnect(context.Background())
		}
	}

	tx.Init(masterDb)
	data = &Data{masterDb: masterDb, slaveDb: slaveDb, redisDb: redisClient,
		rabbitMQChn: rabbitMQChn, kfkProducer: kfkProducer, kfkConsumer: kfkConsumer, mgoDb: mgoDb}
	return data, cleanup, nil
}

func newMasterDB(master *setting.Master) *sqlx.DB {
	var err error
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", master.User, master.Password, master.Host, master.Port, master.DB)
	masterDb, err := sqlx.Connect(master.DriverName, dsn)
	if err != nil {
		panic(err)
	}
	masterDb.SetMaxOpenConns(master.MaxOpenConns)
	masterDb.SetMaxIdleConns(master.MaxIdleConns)
	if err = masterDb.Ping(); err != nil {
		panic(err)
	}
	zap.L().Info("DB init success...")
	return masterDb
}

func newSlaveDB(slave *setting.Slave) []*sqlx.DB {
	count := slave.Count
	var slaveDb []*sqlx.DB
	if count > 0 {
		slaveDb = make([]*sqlx.DB, count)
		var err error
		for i := 0; i < count; i++ {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", slave.Users[i], slave.Passwords[i], slave.Hosts[i], slave.Ports[i], slave.DBs[i])
			slaveDb[i], err = sqlx.Connect(slave.DriverName, dsn)
			if err != nil {
				slaveDb[i].Ping()
				panic(err)
			}
			slaveDb[i].SetMaxOpenConns(slave.MaxOpenConns)
			slaveDb[i].SetMaxIdleConns(slave.MaxIdleConns)
		}
	}
	return slaveDb
}

func newRedis(r *setting.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
	})
	if err := rdb.Ping().Err(); err != nil {
		panic(err)
	}
	zap.L().Info("Redis init success...")
	return rdb
}

func newRabbitMq(r *setting.RabbitMQ) *amqp.Channel {
	if !r.Enabled {
		return nil
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	chn, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	zap.L().Info("RabbitMQ init success...")
	return chn
}

func newKafka(r *setting.Kafka) (sarama.SyncProducer, sarama.ConsumerGroup) {
	if !r.Enabled {
		return nil, nil
	}
	mqConfig := sarama.NewConfig()
	mqConfig.Producer.RequiredAcks = sarama.WaitForLocal
	mqConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	mqConfig.Producer.Return.Successes = true

	clusterCfg := sarama.NewConfig()
	clusterCfg.Consumer.Return.Errors = true
	clusterCfg.Consumer.Offsets.AutoCommit.Enable = true
	clusterCfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	clusterCfg.Consumer.Offsets.Retry.Max = 3

	producer, err := sarama.NewSyncProducer([]string{r.Addrs}, mqConfig)
	if err != nil {
		panic(err)
	}

	consumer, err := sarama.NewConsumerGroup([]string{r.Addrs}, "go-group", clusterCfg)
	if err != nil {
		panic(err)
	}

	zap.L().Info("Kafka init success...")
	return producer, consumer
}

func newMongodb(r *setting.Mongodb) *mongo.Database {
	if !r.Enabled {
		return nil
	}
	uri := fmt.Sprintf("mongodb://%s:%d", r.Host, r.Port)
	options := options.Client().ApplyURI(uri).SetConnectTimeout(2 * time.Second).
		SetAuth(options.Credential{Username: r.User, Password: r.Password, AuthSource: r.DBName})
	client, err := mongo.Connect(context.Background(), options)
	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		panic(err)
	}

	zap.L().Info("Mongodb init success...")
	return client.Database(r.DBName)
}
