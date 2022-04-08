create table if not exists events
(
    timestamp timestamp NOT NULL,
    user_id VARCHAR(88) NOT NULL,
    pixel_color varchar(7) NOT NULL,
    x smallint NOT NULL,
    y smallint NOT NULL,
    x1 smallint,
    y1 smallint,
    PRIMARY KEY (timestamp, user_id, pixel_color, x, y) -- needs to include all fields because of cheaters(?)
);

create temp table tmp_table as select * from events with no data;