CREATE DATABASE IF NOT EXISTS tiktok_chat;

USE tiktok_chat;

CREATE TABLE IF NOT EXISTS messages(
    id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    chat VARCHAR(100) NOT NULL,
    sender VARCHAR(255) NOT NULL,
    text VARCHAR(2000) NOT NULL,
    send_time BIGINT NOT NULL,
    INDEX chat_index (chat),
    INDEX send_time_index (send_time)
    );