import requests
import json
from tqdm import tqdm
from fire import Fire
import random
import multiprocessing


def cast_vote(proto, host, ballot, vote_code, user):
    try:
        requests.get(f"{proto}://{host}/api/election/sse", timeout=0.2)
    except:
        pass
    vote = ballot.copy()
    random.shuffle(vote)
    vote = vote[:-random.randrange(0, len(vote))]
    r = requests.post(f"{proto}://{host}/api/election/current/vote", cookies={"vot-tok": user[2]}, json={"code": vote_code, "vote": vote})
    r.raise_for_status()


def mock_election(admin_user: str, admin_password: str, vote_code: str, n: int=1, host: str="127.0.0.1:9090", proto: str="http", start_manually: bool=False):
    with open("users.json") as f:
        users = json.load(f)

    # login admin
    r = requests.post(f"{proto}://{host}/auth/login/do", data={"studentid": admin_user, "password": admin_password}, allow_redirects=False)
    r.raise_for_status()
    admin_token = r.cookies["vot-tok"]

    for i in range(n):
        print(f"Election {i}")

        # create an election
        r = requests.post(f"{proto}://{host}/api/admin/election", cookies={"vot-tok": admin_token}, json={"roleName": "Bananas"})
        r.raise_for_status()
        election_info = r.json()

        if start_manually:
            input("Please start the election and press enter")
        else:
            # start election
            r = requests.post(f"{proto}://{host}/api/admin/election/start", cookies={"vot-tok": admin_token}, json={"id": election_info["id"], "extraNames": ["Volume Knob", "Power Switch", "Blue Shaver", "Fuzzy Hedgehog"]})
            r.raise_for_status()

        # get ballot
        r = requests.get(f"{proto}://{host}/api/election/current", cookies={"vot-tok": admin_token})
        r.raise_for_status()
        ballot_info = r.json()["ballot"]
        ballot_ids = list(map(lambda x: x["id"], ballot_info))

        # dispatch votes
        with multiprocessing.Pool(processes=8) as pool:
            res = [pool.apply_async(cast_vote, (proto, host, ballot_ids, vote_code, user)) for user in users]
            for x in tqdm(res):
                x.get()

        # end election
        r = requests.post(f"{proto}://{host}/api/admin/election/stop", cookies={"vot-tok": admin_token}, json={"id": election_info["id"]})
        r.raise_for_status()
        results = r.json()
        with open(f"result.{i}.txt", "w") as f:
            f.write(results["result"])


if __name__ == "__main__":
    Fire(mock_election)