CREATE TABLE donations (
   id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
   name VARCHAR(255) NOT NULL,
   photo VARCHAR(255) NOT NULL,
   location VARCHAR(255) NOT NULL,
   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
   userid INT NOT NULL,
   FOREIGN KEY (userid) REFERENCES users(id)
);