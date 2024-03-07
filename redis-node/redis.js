import { createClient } from "redis";
import { QUEUE_NAME } from "./constants.js";
import { writeMessage } from "./rabbitMQ.js";

const client = createClient();

export async function getCounter() {
  return new Promise((resolve, reject) => {
    client.get("COUNTER", (err, res) => {
      if (err) {
        reject(err);
        return;
      }
      res = isNaN(res) ? 0 : res;
      resolve(parseInt(res));
    });
  });
}

export async function incrCounter() {
  client.incr("COUNTER", (err, res) => {
    if (err) {
      throw err;
    }
    // console.log("Value", res);
    enqueue(res);
    writeMessage(res);
    if (!res) {
      throw Error("Error in function incrCounter. Got null result value");
    }
  });

  // let value = (await getCounter()) + 1;
  // // console.log("Value", value);
  // client.set("COUNTER", value, (err, res) => {
  //   if (err) {
  //     throw err;
  //   }
  //   enqueue(res + " " + value);
  //   writeMessage(res + " " + value);
  //   if (!res) {
  //     throw Error("Error in function incrCounter. Got null result value");
  //   }
  // });
}

export async function enqueue(number) {
  // console.log("enqueue number", number);
  client.rpush(QUEUE_NAME, number, (err) => {
    if (err) {
      throw err;
    }
  });
}

export async function dequeue() {
  client.lpop(QUEUE_NAME, (err, result) => {
    if (err) {
      throw err;
    }
  });
}
