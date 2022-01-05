# notifications
Project named “Barley” because its euphony is “大卖”. 

## project description
- provides customer a way to view their purchase history quickly, instead of checking it by looking back our chat history one by one. 
- allows customers to track the package easily.

## features
admin
- admin login
- admin creates user account
- admin create shipping information (general product description)
- admin can trigger shipping notification by sms (copy msg by one-click)

user
- user login
- user check shipping history

## plans
### frontend user
- user authorization (http://example.com/?code=jwt)
- list shipping history page filter, sort 
- list shipping history page - action button - redirect to tracking website

### frontend admin
- admin login 
- admin creates user account
- list shipping history page filter, sort 
- list shipping history page - action button

### database
- design schema
- 
- shipment: 
```
DROP TABLE IF EXISTS `shipment`;
CREATE TABLE `shipment` (
    `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
    `UserID` int(6) NOT NULL,
    `Description` varchar(64),
    `Tracking` varchar(16) NOT NULL,
    `Comment` varchar(255),
    `Date` varchar(16) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
INSERT INTO shipment (id, UserID, Description, Tracking, `Comment`, Date) VALUES(1, 111, "purchased 1* skin care on 12/01", "Udex111111", "Carried by Udex in Canada, will be transferred to Yunda in China", "2021-12-01");
INSERT INTO shipment (id, UserID, Description, Tracking, `Comment`, Date) VALUES(2, 222, "purchased 2* skin care on 12/02", "Udex222222", "Carried by Udex in Canada, will be transferred to Yunda in China", "2021-12-02");
INSERT INTO shipment (id, UserID, Description, Tracking, `Comment`, Date) VALUES(3, 333, "purchased 3* skin care on 12/03", "Udex333333", "Carried by Udex in Canada, will be transferred to Yunda in China", "2021-12-03");
INSERT INTO shipment (id, UserID, Description, Tracking, `Comment`, Date) VALUES(4, 444, "purchased 4* skin care on 12/04", "Udex444444", "Carried by Udex in Canada, will be transferred to Yunda in China", "2021-12-04");
INSERT INTO shipment (id, UserID, Description, Tracking, `Comment`, Date) VALUES(5, 555, "purchased 5* skin care on 12/05", "Udex555555", "Carried by Udex in Canada, will be transferred to Yunda in China", "2021-12-05");
```
- user:
```
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
    `WeChatID` int(6) NOT NULL,
    `WeChatName` varchar(255),
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```
- Admin: (We don't want to make a complicated authorization, so I didn't encode the password)
```
DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin` (
    `id` int(6) unsigned NOT NULL AUTO_INCREMENT,
    `Name` varchar(16) NOT NULL,
    `password` varchar(16) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;
```

### api
#### frontend
- GET /api/v1/users/me
- GET /api/v1/shippings
#### admin
- POST /api/admin/v1/users 
- POST /api/admin/v1/shippings 
- POST /api/v1/shippings/actions/resend-notification (for admin view)
- GET /api/admin/v1/shippings
#### I have changed them to:
#### user
- to be done
#### admin
- POST /api/v1/admin/shipment (to create a new shipment)
- GET /api/v1/admin/shipment (to view all shipments)
- GET /api/v1/admin/shipments/{id} (to view one shipment)
- DELETE /api/v1/admin/shipments/{id} (to delete one shipment)
- PATCH /api/v1/admin/shipments/{id} (to update one shipment)
- GET /api/v1/admin/users (to view all users)
- POST /api/v1/admin/user (to create a new user)
- to be done: POST /api/v1/admin/login (to login to admin account)

### lib
#### ant design pro
- https://pro.ant.design/
#### rest api
- https://github.com/gorilla/mux
- example
- https://medium.com/the-andela-way/build-a-restful-json-api-with-golang-85a83420c9da
- https://hugo-johnsson.medium.com/rest-api-with-golang-mux-mysql-c5915347fa5b
#### mysql golang driver 
- https://github.com/go-sql-driver/mysql
- example
- https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html
#### convert payload to json format
- https://www.sohamkamani.com/golang/json/
#### rest api documentation
- https://editor.swagger.io
- https://swagger.io/docs/specification/2-0/basic-structure/
- https://medium.com/@amirm.lavasani/restful-apis-tutorial-of-openapi-specification-eeada0e3901d
#### docker
- https://labs.play-with-docker.com
- https://www.docker.com/101-tutorial
- https://hub.docker.com/repository/docker/gentledonkey/101-todo-app
