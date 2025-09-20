import http from 'k6/http'
import {sleep} from 'k6'

export const options = {
    iterations: 10_000,   
    vus: 10,
    batch: 5_000,
}

export default function () {
    http.get("http://sample-service.sample:5000/")
    sleep(0.1)
}
