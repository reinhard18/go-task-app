# Register a new user

curl -X POST http://localhost:8080/api/register \
 -H "Content-Type: application/json" \
 -d '{
"username": "newuser",
"password": "newpassword123"
}'

# The response will include a JWT token and username if successful

# You can then try to login with the new credentials

curl -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/json" \
 -d '{
"username": "newuser",
"password": "newpassword123"
}'

curl -X POST http://localhost:8080/tasks \
-H "Authorization: Bearer YOUR_TOKEN" \
-H "Content-Type: application/json" \
-d '{"title":"Test Task","description":"This is a test task","status":"pending"}'

curl http://localhost:8080/tasks \
-H "Authorization: Bearer YOUR_TOKEN"

curl -X PUT http://localhost:8080/task/1 \
-H "Authorization: Bearer YOUR_TOKEN" \
-H "Content-Type: application/json" \
-d '{"title":"Updated Task","description":"This is updated","status":"completed"}'

curl -X DELETE http://localhost:8080/task/1 \
-H "Authorization: Bearer YOUR_TOKEN"
