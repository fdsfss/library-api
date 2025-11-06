CREATE TABLE authors(
                        ID             UUID PRIMARY KEY,
                        full_name       TEXT NOT NULL CHECK ( full_name <> '' ),
                        nick_name       TEXT NOT NULL CHECK ( nick_name <> '' ),
                        specialization TEXT NOT NULL CHECK ( specialization <> '' ));

INSERT INTO authors (ID, full_name, nick_name, specialization) VALUES
                                                                   ('4ce0ddc1-ed52-4173-8e82-e32926ddff2e', 'Stephen King', 'The King of Horror', 'Horror Fiction'),
                                                                   ('4656a695-fb18-4929-b37c-a6c32a5280e8', 'J.K. Rowling', 'Jo', 'Fantasy Fiction'),
                                                                   ('d23bdad0-0d90-47b2-b202-8fa6eea08c80', 'Isaac Asimov', 'The Good Doctor', 'Science Fiction'),
                                                                   ('d5345aed-3533-4ea2-a94e-50722589d956', 'Agatha Christie', 'Queen of Mystery', 'Mystery Fiction'),
                                                                   ('37a1cf27-51c6-4a04-9018-6b05e6a1dd46', 'Ernest Hemingway', 'Papa', 'Literary Fiction');





