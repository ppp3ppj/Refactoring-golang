INSERT INTO "Person" (key, name, description, image, traits, tags)
VALUES
    (
        'minji',
        'Minji',
        'She’s the cleanup queen of NewJeans. Her nickname is Teddy Bear.',
        'path/to/minji.png',
        '[{"Personality": "ESTJ", "Like": "eating ice cream", "Zodiac Sign": "Taurus", "Emoji": "🐻", "Color": "Blue" }]',
        '{"minji", "newjeans"}'
    ),
    (
        'hanni',
        'Hanni',
        'She’s good at sleeping fast anywhere, even while sitting. Her nickname is Pigtails.',
        'path/to/hanni.png',
        '[{"Personality": "INFP", "Like": "wearing hoodies", "Zodiac Sign": "Libra", "Emoji": "🐰", "Color": "Pink" }]',
        '{"hanni", "newjeans"}'
    ),
    (
        'danielle',
        'Danielle',
        'She wants to start a surfing club with Hanni since she loves surfing so much. Her nickname is Dana, Danette, Dani, Dania, Danita, Danka, Danna, Danni, Dannie, Danny, Dany, Daša, Dell, Della, El, Ellie, Elly, Nelia, Nella.',
        'path/to/danielle.png',
        '[{"Personality": "ENFP", "Like": "drawing, painting, listening to music, swimming, and talking with the members", "Zodiac Sign": "Aries", "Emoji": "🐶", "Color": "Yellow" }]',
        '{"danielle", "newjeans"}'
    ),
    (
        'haerin',
        'Haerin',
        'She thinks she’s very unpredictable. Her nickname is Kitty Kang.',
        'path/to/haerin.png',
        '[{"Personality": "INTP", "Like": "listening to music and reading", "Zodiac Sign": "Taurus", "Emoji": "🐱", "Color": "Green" }]',
        '{"haerin", "newjeans", "cat"}'
    ),
    (
        'hyein',
        'Hyein',
        'She loves Harry Potter, and has a bunch of its books in English and Korean. Her nickname is Hyeni.',
        'path/to/hyein.png',
        '[{"Personality": "ISFP", "Like": "talking walks, taking photos of mostly the sky and members, and looking up movies", "Zodiac Sign": "Taurus", "Emoji": "🐣", "Color": "Purple" }]',
        '{"hyein", "newjeans"}'
    )
ON CONFLICT (key) DO UPDATE
SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    image = EXCLUDED.image,
    traits = EXCLUDED.traits,
    tags = EXCLUDED.tags;
