import http from "k6/http";
import {check, group} from "k6";
import {randomString} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";

export const BASE_URL = "http://localhost:8080/api/v1";

export const requestConfigWithTag = (tag) => ({
    headers: {
        Content_Type: "application/json",
    },

    tags: Object.assign(
        {name: tag},
    ),
});

function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min)
}

export function createSong() {
    const res = http.post(`${BASE_URL}/songs`,
        JSON.stringify({
            author: randomString(8),
            duration: `${randomInt(1, 5)}:${randomInt(0, 59)}`,
            release: `${randomInt(2000, 2024)}`,
            title: randomString(8),
        }),
        requestConfigWithTag("createSong"),
    );

    const ok = check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
        "id is present": (r) => r.json("id") !== "",
    });

    if (!ok) {
        console.log(`creating a song failed: ${res.status} ${res.body}`);
        return ""
    }

    return res.json("id");
}

export function getAllSongs() {
    const res = http.get(`${BASE_URL}/songs`, requestConfigWithTag("getAllSongs"));

    const ok = check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
    });

    if (!ok) {
        console.log(`get all songs failed: ${res.status} ${res.body}`);
    }
}

export function getSong(songId) {
    const res = http.get(`${BASE_URL}/songs/${songId}`, requestConfigWithTag("getSong"));

    const ok = check(res, {
        "status is 200": (r) => r.status === 200,
        "response is in json": (r) => r.headers["Content-Type"] === "application/json; charset=utf-8",
        "validate id": (r) => r.json("id") === songId,
    });

    if (!ok) {
        console.log(`get song failed: ${res.status} ${res.body}`);
    }
}

export function deleteSong(songId) {
    const res = http.del(`${BASE_URL}/songs/${songId}`, requestConfigWithTag("deleteSong"));

    const ok = check(res, {
        "status is 204": (r) => r.status === 204,
    });

    if (!ok) {
        console.log(`delete song failed: ${res.status} ${res.body}`);
    }
}

export const options = {
    scenarios: {
        stress_test: {
            executor: "ramping-arrival-rate",
            stages: [
                {target: 100, duration: "10s"},
                {target: 200, duration: "20s"},
                {target: 400, duration: "1m"},
                {target: 300, duration: "1m"},
                {target: 300, duration: "1m"},
                {target: 500, duration: "2m"},
                {target: 250, duration: "30s"},
                {target: 100, duration: "2m"},
                {target: 0, duration: "30s"},
            ],
            preAllocatedVUs: 230,
        },
    },
};

export default () => {
    let songIds = [];

    group("01. Create songs", () => {
        const id = createSong();
        if (id !== "") {
            songIds.push(id)
        }
    });

    group("02. Get all songs", () => {
        getAllSongs(songIds);
    });

    group("03. Get song", () => {
        getSong(songIds[randomInt(0, songIds.length - 1)]);
    });

    group("04. Delete song", () => {
        deleteSong(songIds[randomInt(0, songIds.length - 1)]);
    });
}
