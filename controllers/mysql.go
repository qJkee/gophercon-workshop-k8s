package controllers

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/workshop/mysql-operator/api/v1alpha1"

	_ "github.com/go-sql-driver/mysql"
)

func (r *CustomMysqlReconciler) manageReplication(cr *v1alpha1.CustomMysql) error {
	conf := mysql.NewConfig()
	conf.User = "root"
	conf.Passwd = "root_password"
	conf.Net = "tcp"
	// ТУТ НУЖНО ПРАВИЛЬНО ЗАПИСАТЬ АДРЕС, HINT ЕСТЬ =)
	// conf.Addr = "$(Statefulset_NAME)-$(POD_NUM).some-name-mysql.$(NAMESPACE_NAME).svc.cluster.local:33061"
	conf.Params = map[string]string{"interpolateParams": "true"}

	mysqlDB, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return errors.Wrap(err, "cannot connect to host")
	}
	return nil
}

// TODO LIST
// 1. Подключится к mysql ноде 0 и выполнить эти комманды

// SET GLOBAL group_replication_group_name='aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
// SET GLOBAL group_replication_local_address='$(Statefulset_NAME)-0.some-name-mysql.$(NAMESPACE_NAME).svc.cluster.local:33061';
// SET GLOBAL group_replication_group_seeds='$(Statefulset_NAME)-0.some-name-mysql.$(NAMESPACE_NAME).svc.cluster.local:33061';
// SET GLOBAL group_replication_ip_allowlist='0.0.0.0/0';
// SET GLOBAL group_replication_recovery_get_public_key=ON;
// SET GLOBAL group_replication_single_primary_mode=OFF;
// SET GLOBAL group_replication_bootstrap_group=ON;
// START GROUP_REPLICATION USER='root', PASSWORD='root_password';
// SET GLOBAL group_replication_bootstrap_group=OFF;

// 2. ПОДКЛЮЧИТСЯ К ОСТАЛЬНЫМ НОДАМ И ВЫПОЛНИТЬ
// SET GLOBAL group_replication_group_name='aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
// SET GLOBAL group_replication_local_address='$(Statefulset_NAME)-1.some-name-mysql.$(NAMESPACE_NAME).svc.cluster.local:33061';
// SET GLOBAL group_replication_group_seeds='$(Statefulset_NAME)-0.some-name-mysql.$(NAMESPACE_NAME).svc.cluster.local:33061,$(Statefulset_NAME)-1.some-name-mysql.$(NAMESPACE_NAME).svc.cluster.local:33061';
// SET GLOBAL group_replication_ip_allowlist='0.0.0.0/0';
// SET GLOBAL group_replication_recovery_get_public_key=ON;
// SET GLOBAL group_replication_single_primary_mode=OFF;
// CHANGE MASTER TO MASTER_USER='root', MASTER_PASSWORD='root_password' FOR CHANNEL 'group_replication_recovery';
// RESET MASTER;
// START GROUP_REPLICATION USER='root', PASSWORD='root_password';

// 3. ПРОВЕРИТЬ)
// kubectl exec -it some-name-mysql-0 -- sh -c "mysql -uroot -proot_password -h127.0.0.1"
// CREATE DATABASE test;
// CREATE TABLE test.t1 (c1 INT NOT NULL PRIMARY KEY) ENGINE=InnoDB;
// INSERT INTO test.t1 VALUES (1);
// SELECT * FROM test.t1;
// ​
// 3.1 Этот селект должен вернуть значения, которые мы записали в поду 0
// kubectl exec -it some-name-mysql-1 -- sh -c "mysql -uroot -proot_password -h127.0.0.1"
// SELECT * FROM test.t1;
// INSERT INTO test.t1 VALUES (2);
// SELECT * FROM test.t1;

// 4. Как проверить что репликация уже настроена и не нужно делать это все еще раз?
// select * from performance_schema.replication_group_members;
// Этот запрос вернет мемберов репликации, отсюда и танцуем)

// 5. OVERTASK
// Добавить возможность указывать N количество реплик(сейчас мы указываем лишь 1 мастер и 1 реплику)
// Т.е выполнить правильно пункт 2 для всех подов 1-cr.Size
