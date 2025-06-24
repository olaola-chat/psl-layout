#!/bin/bash

# TODO -tables 这一行要替换成自己的表名

gentool -dsn "root:123456@tcp(127.0.0.1:3306)/xianshi" \
  -tables "xs_user_coupon,xs_user_coupon_log" \
  -outPath internal/data/layout/gen \
  -fieldSignable
