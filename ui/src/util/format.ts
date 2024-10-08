import { default as dayjs } from "dayjs";

export const formatDate = (date: number) =>
  dayjs(date).format("MMMM D, h:mm A");

export const getDaysInMonth = (year: number, month: number) => {
  const daysInMonth = new Date(year, month + 1, 0).getDate();
  return Array.from(
    { length: daysInMonth },
    (_, i) => new Date(year, month, i + 1)
  );
};
