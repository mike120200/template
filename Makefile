
# 生成pg数据库
docker:
	docker run -it --name calorie_postgres_backup \
  	-d \
 	--restart=always \
  	-p 5555:5432 \
  	--network qnear \
  	-e POSTGRES_USER=calorie \
  	-e POSTGRES_DB=calorie_db \
  	-e POSTGRES_PASSWORD=calorie_2025 \
  	postgres:latest

#启动redis
redis:
	docker run -d --name calorie_redis \
  	-p 6379:6379 \
  	--restart=always \
  	-e REDIS_PASSWORD=calorie_2025 \
  	redis:latest redis-server --requirepass calorie_2025


# 下载依赖
deps:
	go mod tidy

# 运行
run:
	go run main.go dev

# 打包
build:
	go build -o calorie

test:
	go test -parallel 1 ./...

sc:
	go run script/gen_model.go