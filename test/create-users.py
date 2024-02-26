import string
from tqdm import tqdm
import json
import requests
from fire import Fire
import random

def create_users(n: int=100, host: str="127.0.0.1:9090"):
    res = []

    url = f"http://{host}/auth/login/do"

    for i in tqdm(range(n)):
        sid = str(random.randint(1000000, 3000000))
        pwd = "".join(random.choices(string.ascii_letters + string.punctuation, k=10))
        r = requests.post(url, data={
            "studentid": sid,
            "password": pwd,
            "passwordconf": pwd,
            "fname": "Martin",
            "lname": "Martinson",
        }, allow_redirects=False)
        r.raise_for_status()

        res.append([sid, pwd, r.cookies["vot-tok"]])

    with open("users.json", "w") as f:
        json.dump(res, f, indent=2)

if __name__ == "__main__":
    Fire(create_users)