import { init, writeMessage } from "./rabbitMQ.js";
import { THREAD_COUNT, QUEUE_NAME } from "./constants.js";
import { enqueue, dequeue, incrCounter } from "./redis.js";

async function main() {
  for (let i = 0; i < THREAD_COUNT; i++) {
    initThread(i);
  }
}

async function initThread(threadNum) {
  console.log(`Thread ${threadNum} has started`);
  const channel = await init();

  setInterval(async () => {
    await incrCounter();

    channel.consume(
      QUEUE_NAME,
      async (msg) => {
        const number = msg.content.toString();

        console.log(
          `Got ${number} from ${QUEUE_NAME} with thread ${threadNum}`
        );

        // await enqueue(number);

        // const popNumber = await dequeue();

        // if (popNumber) {
        //   console.log(`Popped ${popNumber} from ${QUEUE_NAME}`);
        //   writeMessage(popNumber);
        // }
      },
      { noAck: true }
    );
  }, 2000);
}

main();
