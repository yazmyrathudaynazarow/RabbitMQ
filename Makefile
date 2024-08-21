# Runs a Docker container with RabbitMQ server including the management UI.
rabbitmq-server:
	sudo docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management

rabbitmq:
	sudo docker start rabbitmq

# user percy. password secret
add-user-percy:
	sudo docker exec rabbitmq rabbitmqctl add_user percy secret

percy-admin:
	sudo docker exec rabbitmq rabbitmqctl set_user_tags percy administrator

add-vhost-customers:
	sudo docker exec rabbitmq rabbitmqctl add_vhost customers

set-permissions-percy:
	sudo docker exec rabbitmq rabbitmqctl set_permissions -p customers percy ".*" ".*" ".*"

declare-exchange:
	sudo docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=topic -u percy -p secret durable=true

permission_exchange:
	sudo docker exec rabbitmq rabbitmqctl set_topic_permissions -p customers percy customer_event "^customer.*" "^customer.*"

# type = fanout
create-exchange:
	sudo docker exec rabbitmq rabbitmqadmin declare exchange --vhost=customers name=customer_events type=fanout -u percy -p secret durable=true