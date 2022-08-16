CREATE TABLE IF NOT EXISTS cars(
    id SERIAL NOT NULL PRIMARY KEY,
    make VARCHAR NOT NULL,
    model VARCHAR NOT NULL,
    package VARCHAR NOT NULL,
    color VARCHAR NOT NULL,
    identification VARCHAR NOT NULL,
    year INTEGER NOT NULL,
    category VARCHAR NOT NULL,
    mileage INTEGER NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP                             
);
