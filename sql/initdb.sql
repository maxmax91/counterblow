CREATE TABLE IF NOT EXISTS rules (
  rule_id INT NOT NULL GENERATED ALWAYS AS IDENTITY,
  rule_type INT NOT NULL, -- current implementations 1: round robin, 2: ip hash -- maybe is better string or enum?
  rule_ipaddr varchar(16) NULL, -- null = no filters
  rule_subnetmask int NULL, -- null = no filter
  rule_servers varchar NOT NULL,
  rule_source varchar NULL DEFAULT NULL, -- catch the request only if match this regex. Default .*
  rule_dest varchar NULL DEFAULT NULL, -- rewrite the request using this template. Default $0
  PRIMARY KEY (rule_id)
);


CREATE TABLE IF NOT EXISTS hits (
  hit_from varchar(250) NOT NULL, -- requesting ip address
  hit_to varchar(250) NOT NULL, -- routing address
  hit_datetime timestamp without time zone default (now() at time zone 'utc') -- utc time so there is no confusion
);

-- test rules
INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES (1, '0.0.0.0', 0, 'microsoft.it:80', '/test1/(.*)', '$1/rewrote/');
INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES (1, '0.0.0.0', 0, 'google.it:80', '/test2/(.*)', '$1/rewrote/');
INSERT INTO rules (rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) VALUES (1, '0.0.0.0', 0, 'google.it:80,microsoft.it:80,tesla.com:80', '.*', '$0');


-- 
-- localhost:8080/test1/(something) will redirect to microsoft.it:80/(something)/rewrote/
-- localhost:8080/test2/(something) will redirect to google.it:80/(something)/rewrote/
-- everything else will redirect balancing to google.it,microsoft.it and tesla.com