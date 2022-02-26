import random
import json
import string

NUM = 10000
STRING_LEN = 30
NUM_MAPS = 1000


def get_random_string():
    letters = string.ascii_letters
    return ''.join(random.choice(letters) for _ in range(STRING_LEN))


def get_random_maps():
    return {get_random_string(): random.randrange(100000) for _ in range(NUM_MAPS)}


out = {'tests': []}

if __name__ == '__main__':
    for i in range(NUM):
        out['tests'].append(
            {
                'id': random.randrange(100000),
                'frac': random.random(),
                'name': get_random_string(),
                'maps': get_random_maps()
            }
        )

    with open('json_init.json', 'w') as f:
        json.dump(out, f)
