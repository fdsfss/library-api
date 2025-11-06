CREATE TABLE members(
                        ID             UUID PRIMARY KEY,
                        full_name       TEXT NOT NULL CHECK ( full_name <> '' ));

INSERT INTO members (ID, full_name) VALUES
                                        ('5d574a92-4b78-46eb-8ab0-02709b710b15', 'John Smith'),
                                        ('56013726-9dd0-436a-8722-f5e5a9896dc6', 'Emily Johnson'),
                                        ('0d3a7572-eaa6-4e74-904f-30c0f2842981', 'Michael Brown'),
                                        ('2b3692bc-07f0-4748-9847-13dc26409066', 'Sarah Williams'),
                                        ('80eb37ea-d1ba-4906-bd12-1f84a29b22a0', 'David Martinez');