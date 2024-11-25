CREATE TABLE IF NOT EXISTS songs (
    id SERIAL,
    "group" TEXT NOT NULL,
    song TEXT PRIMARY KEY NOT NULL UNIQUE ,
    releasedate TEXT NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL
);