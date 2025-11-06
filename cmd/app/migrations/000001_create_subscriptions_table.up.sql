CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INT NOT NULL,
    user_id UUID NOT NULL,
    start_date CHAR(7) NOT NULL, 
    end_date CHAR(7) NOT NULL    
);