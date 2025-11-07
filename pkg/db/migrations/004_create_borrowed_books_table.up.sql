CREATE TABLE borrowed_books(
                               member_id UUID REFERENCES members(ID),
                               book_id UUID REFERENCES books(ID));

INSERT INTO borrowed_books (member_id, book_id) VALUES
                                                    ((SELECT ID FROM members WHERE full_name = 'John Smith'), (SELECT ID FROM books WHERE title = 'The Shining')),
                                                    ((SELECT ID FROM members WHERE full_name = 'John Smith'), (SELECT ID FROM books WHERE title = 'IT')),
                                                    ((SELECT ID FROM members WHERE full_name = 'John Smith'), (SELECT ID FROM books WHERE title = 'Foundation')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Emily Johnson'), (SELECT ID FROM books WHERE title = 'Harry Potter and the Philosopher''s Stone')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Emily Johnson'), (SELECT ID FROM books WHERE title = 'Harry Potter and the Chamber of Secrets')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Emily Johnson'), (SELECT ID FROM books WHERE title = 'Murder on the Orient Express')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Michael Brown'), (SELECT ID FROM books WHERE title = 'The Old Man and the Sea')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Michael Brown'), (SELECT ID FROM books WHERE title = 'A Farewell to Arms')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Sarah Williams'), (SELECT ID FROM books WHERE title = 'I, Robot')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Sarah Williams'), (SELECT ID FROM books WHERE title = 'And Then There Were None')),
                                                    ((SELECT ID FROM members WHERE full_name = 'Sarah Williams'), (SELECT ID FROM books WHERE title = 'The Sun Also Rises')),
                                                    ((SELECT ID FROM members WHERE full_name = 'David Martinez'), (SELECT ID FROM books WHERE title = 'Death on the Nile')),
                                                    ((SELECT ID FROM members WHERE full_name = 'David Martinez'), (SELECT ID FROM books WHERE title = 'Carrie'));