create table if not exists events
(
    timestamp timestamp NOT NULL,
    user_id VARCHAR(88) NOT NULL,
    pixel_color varchar(7) NOT NULL,
    -- moderator actions will be added as many discrete pixels in the same timestamp
    x smallint NOT NULL,
    y smallint NOT NULL,
    PRIMARY KEY (timestamp, user_id, pixel_color, x, y) -- needs to include all fields because of cheaters(?)
);
-- todo flag moderators in some way