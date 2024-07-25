INSERT INTO "Person" (key, name, description, image, traits, tags)
VALUES
    (
        'newjean',
        'Newjean',
        'An up-and-coming artist known for her innovative music and style.',
        'path/to/newjean_image.png',
        '["talented", "innovative"]',
        '{"artist", "musician", "fashion-icon"}'
    )
ON CONFLICT (key) DO UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    image = EXCLUDED.image,
    traits = EXCLUDED.traits,
    tags = EXCLUDED.tags;

INSERT INTO "Person" (key, name, description, image, traits, tags)
VALUES
    (
        'alice',
        'Alice',
        'A talented painter with a passion for abstract art.',
        'path/to/alice_image.png',
        '["creative", "abstract"]',
        '{"painter", "artist"}'
    ),
    (
        'bob',
        'Bob',
        'A seasoned chef known for his expertise in Italian cuisine.',
        'path/to/bob_image.png',
        '["skilled", "cuisine-specialist"]',
        '{"chef", "italian-cuisine"}'
    )
ON CONFLICT (key) DO UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    image = EXCLUDED.image,
    traits = EXCLUDED.traits,
    tags = EXCLUDED.tags;
