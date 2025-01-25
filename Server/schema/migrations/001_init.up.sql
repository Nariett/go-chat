CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name CHARACTER VARYING(30) UNIQUE,
                       password CHARACTER VARYING(30)
);

CREATE TABLE activity(
                         id SERIAL PRIMARY KEY,
                         idUser INTEGER NOT NULL,
                         date TIMESTAMP NOT NULL,
                         FOREIGN KEY (idUser) REFERENCES users (id)
);

CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          sender_id INT NOT NULL,
                          recipient_id INT NOT NULL,
                          content TEXT NOT NULL,
                          sent_at TIMESTAMP NOT NULL,
                          read_at TIMESTAMP,
                          FOREIGN KEY (sender_id) REFERENCES users (id),
                          FOREIGN KEY (recipient_id) REFERENCES users (id)
);

INSERT INTO users (name, password) VALUES
                                       ('Саша', '123456'),
                                       ('Олег', '1111'),
                                       ('Коля', '4321');

INSERT INTO activity(idUser, date)VALUES
                                      (1 , NOW()),
                                      (2 , NOW()),
                                      (3 , NOW());