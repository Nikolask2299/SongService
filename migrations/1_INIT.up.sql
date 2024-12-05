CREATE TABLE IF NOT EXISTS groups (
    Id SERIAL,
    "group" TEXT,
    PRIMARY KEY ("group")
);


CREATE TABLE IF NOT EXISTS songs (
    Id SERIAL,
    "group" TEXT,
    song TEXT PRIMARY KEY UNIQUE ,
    releasedate DATE,
    text TEXT,
    link TEXT,
    FOREIGN KEY ("group") REFERENCES groups("group") 
);

