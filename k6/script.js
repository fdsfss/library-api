import http from 'k6/http';
import { Trend } from "k6/metrics";
import { check, sleep } from 'k6';

let myTrend = new Trend("my_trend");

export let options = {
    stages: [
        { duration: '1m', target: 100 },
    ],
};

function generateString() {
    let result           = '';
    let characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz ';
    let charactersLength = characters.length;
    for ( let i = 0; i < 10; i++ ) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }

    return result;
}

function generateQueryBody() {
    let FullName = generateString();
    let NickName = generateString();
    let Specialization = generateString();

    let query = {
        "full_name": FullName,
        "nick_name": NickName,
        "specialization": Specialization
    };

    return query
}

let headers = {
    'Content-Type': 'application/json',
};

export default function () {
    let res = http.post('http://library-api-app:8080/author', JSON.stringify(generateQueryBody()), { headers: headers });
    if (res.status !== 201) {
        console.log(JSON.stringify(res.body));
    }
    myTrend.add(res.timings.sending + res.timings.receiving);
    check(res, { 'status was 201': (r) => r.status == 201 });
    sleep(0.1);
}
