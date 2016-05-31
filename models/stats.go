package models

type UserEmojis struct{
  User            string
  Positive        float32
  Negative        float32
  PositiveCount   int32
  NegativeCount   int32
}
