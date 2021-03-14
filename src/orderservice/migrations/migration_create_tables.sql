USE 'orders';
CREATE TABLE orders (
    order_id VARCHAR(40),
    time INT,
    cost INT);
CREATE TABLE order_item (
    id MEDIUMINT NOT NULL AUTO_INCREMENT,
    menu_id VARCHAR(40),
    quantity INT,
    PRIMARY KEY(id));
CREATE TABLE item_in_order(
    order_id VARCHAR(40),
    item_id MEDIUMINT);