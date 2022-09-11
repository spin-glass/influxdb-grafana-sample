# 1. 利用方法
1. [Eco-Totem Broadway Bicycle Count](https://data.cambridgema.gov/Transportation-Planning/Eco-Totem-Broadway-Bicycle-Count/q8v9-mcfg)からcsvデータを取得し、dataフォルダに配置する
2. ```docker-compose up -d```
3. http://localhost:8086 にアクセスし、docker-compose.ymlに設定してあるUsername, PasswordでInfluxDBにログイン
4. `bike`という名称のBucketを作成する
5. http://localhost:3000 にアクセスし、Username:admin, Password:adminでGrafanaにログイン
6. InfluxDB と接続する（tokenはdocker-compose.yml記載）
7. ```go run src/import_data```
8. InfluxDBやGrafanaのダッシュボードで可視化を行う。
