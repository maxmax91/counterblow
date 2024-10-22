## About

Project built with go+wails+React.

This is a load balancer with the possibility to rewrite the URL and add personalized rules through a friendly interface (and a postgres DB behind).

![CounterBlow Load Balancer](./doc/screenshot1.png)

# Theory and possibilities

- Round robin - redistribuite traffic equally among different servers
- Less latency - redirect to the server with less latency (nearest)
- Maximum throughput - redirect to the server with maximum throughput (less stressed)
- Round robin adaptative - Using all those criteria and a alghorithm eg. 
    Round robin probability = 
    
        (throughput (kb/s)*throughput_const  - latency (ms)*latency_const + round_robin_const)
        -----------------------------------------------------------
        summation(throughput (kb/s)*throughput_const - latency (ms)*latency_const + round_robin_const)

throughput_const and latency_const decide how much thos two variables are considered.

- Source ip hash - to assure that the same source ip goes to the same server.
- Smallest load - The server can send diagnostic packet to the load balancer to nofity the load
                 and the load balancer can on the basis of the load (maybe with the same algorithms as before).

Observations:
- Can be use the round robin adaptative and assure the server does not change between users? How?
    1) can be kept a lookout table between ip and server, with a expiration time
    2) can be opened a session using cookies

# Actually implemented stuff
Implemented:
- Ability to add and remove rules via UI with a supporting postgres db
- Ability to rewrite urls and write different rules with different servers alternating

# How to use it
Build the database and 
        
    docker compose up

Then start the program

## How to debug it
    wails dev

Launch.json is configured to run the developement-built binary to attach the VS Code debugger.

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.


# Test data
There are some test data in the database that are wrote when docker-compose is run for the first time.

    -- test rules
    INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES (1, '0.0.0.0', 0, 'microsoft.it:80', '/test1/(.*)', '$1/rewrote/');
    INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES (1, '0.0.0.0', 0, 'google.it:80', '/test2/(.*)', '$1/rewrote/');
    INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES (1, '0.0.0.0', 0, 'google.it:80,microsoft.it:80,tesla.com:80', '.*', '$0');


    -- 
    -- localhost:8080/test1/(something) will redirect to microsoft.it:80/(something)/rewrote/
    -- localhost:8080/test2/(something) will redirect to google.it:80/(something)/rewrote/
    -- everything else will redirect balancing to google.it,microsoft.it and tesla.com

## Delete test data
Shutdown docker-compose, then 

        rm -rf ./postgres-data
to reset the database.

# Sources & references

- Go tour
https://go.dev/tour/welcome/1

- Proxy in go
https://reintech.io/blog/creating-simple-proxy-server-with-go

- Load balancer in go
https://medium.com/@leonardo5621_66451/building-a-load-balancer-in-go-1c68131dc0ef

- Reverse proxy in go
https://medium.com/trendyol-tech/golang-ile-custom-reverse-proxy-yapmak-7a4198fe86fc#

- Go and proxy servers
https://eli.thegreenplace.net/2022/go-and-proxy-servers-part-1-http-proxies/

- Wails
https://wails.io/docs/introduction/

- Database drivers (Postgres)
https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/

- React documentation
https://react.dev/learn/tutorial-tic-tac-toe


# Package installations

1) Install go
https://go.dev/doc/install
echo "export PATH=\$PATH:/usr/local/go/bin:$HOME/go/bin" >> $HOME/.profile

Try with $go version

2) Install npm https://nodejs.org/en/download/
Try with $npm --version

2) Install wails
This will try to install also gcc build tools, libgtk3 and libwebkit
go install github.com/wailsapp/wails/v2/cmd/wails@latest

Try with $wails version

3) Install packages
Postgres driver

go get -u github.com/lib/pq

4) Starting postgres server with sample rules

docker-compose up

4) Configuration
You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

# TODOs
- Gracefully stop of the server (not yet implemented)
https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve

- Consider to interace the database with a ORM?
https://gorm.io/index.html

- UDP simple server for receiving diagnostic stuff?
https://gist.github.com/miekg/d9bc045c89578f3cc66a214488e68227

- React live chart for visualizing reverse proxy usage
https://apexcharts.com/react-chart-demos/line-charts/realtime/

- Regex search and replace for something similar to mod_rewrite
Done
https://www.geeksforgeeks.org/golang-replacing-all-string-which-matches-with-regular-expression/

