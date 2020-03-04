create database if not exists repopool default character set utf8;
use repopool;
set names utf8;
grant all privileges on repopool.* to 'repopool'@'%' identified by '%$&4567rty$%&*()FG';
CREATE TABLE IF NOT EXISTS syncpool(
    id int UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    harborRepo VARCHAR(80) NOT NULL,
    swrRepo VARCHAR(80) NOT NULL,
    syncStatus TEXT(30) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);
