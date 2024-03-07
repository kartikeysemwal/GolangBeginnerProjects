import amqp from "amqplib";
import { QUEUE_NAME } from "./constants.js";

export async function init() {
  const connection = await amqp.connect("amqp://localhost");
  const channel = await connection.createChannel();

  await channel.assertQueue(QUEUE_NAME, { durable: true });

  return channel;
}

export async function writeMessage(msg) {
  const channel = await init();
  msg = String(msg);
  channel.sendToQueue(QUEUE_NAME, Buffer.from(msg), { persistent: true });
}
