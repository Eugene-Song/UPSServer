import faker from 'faker';

function generateRandomRecord(userId: number) {
  const statuses = ['In Transit', 'Delivered', 'Returned', 'Cancelled'];
  const status = statuses[Math.floor(Math.random() * statuses.length)];
  const trackingNumber = faker.random.alphaNumeric(10);
  const date = faker.date.past(1);
  const targetAddr = faker.address.streetAddress();

  return {
    status,
    tracking_number: trackingNumber,
    date,
    targetaddr: targetAddr,
    user_id: userId,
  };
}
