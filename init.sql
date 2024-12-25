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

INSERT INTO users (name, password) VALUES 
('Саша', '123456'),
('Олег', '1111'),
('Коля', '4321');

INSERT INTO activity(idUser, date)VALUES 
(1 , NOW()),
(2 , NOW()),
(3 , NOW());