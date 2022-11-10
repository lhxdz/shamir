# shamir-tools
this is a tool for shamir (k, n) threshold scheme

# 使用方式
## 加密：
输入门限值t、密钥个数n、秘密 即可加密，将会生成一个必须密钥necessary_key、n个密钥(每个密钥包含x、y)，其中任意t个密钥可以恢复秘密

````
lhx@DESKTOP-0GALLEM:~$ shamir encrypt -t 2 -n 3 "this is a secret.同时可以使用中文。"
necessary key: 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_MV0yI9PnICUSbyxh9iceG0W0HlBn6PDFcfsDDWIHFs3
+-------------------------+-----------------------------------------------------------------------------------------------------------+
|          KEY X          |                                                   KEY Y                                                   |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
| pBE2KvVaYN_569FAGOeEKr  | BVBszEaZArdp0XuVWgxQkDUSr6ylDCgi5t2i7nH4nMEh9ehZeyrPSbs8OnOu_pENXjVJguH45QnxV4k6P5gkxfhm989b23uMqVCd5dWc  |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
| 8PXwCwkVb7L_2huJ6o8FBgE | uvzHiQfto4zzEJvjhKhCd30XIosCKe2XR3Uh8ULKXULqgA9skmQqlQ0yTWHk_vf0gf0Elij0iedGYRKHrbpTg5CnYvItrt1J9NalREqy  |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
| aWTt0gYv09K_8Id7tS4dstb | 1CvtrfulBa0ObWofpE4ygNK4bQdDWetizrXSXAdWbLTC4T92YAlJpLvFvkMdx_iQcV2Pdrfgb6iCbv9zoBDpXfekqW1f5fVzj7gs3WSCp |
+-------------------------+-----------------------------------------------------------------------------------------------------------+
````

## 解密：
分别输入t个密钥x和y，其中第i个x和第i个y是一对密钥，同时输入necessary_key，即可恢复秘密

````
lhx@DESKTOP-0GALLEM:~$ shamir decrypt -n 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_MV0yI9PnICUSbyxh9iceG0W0HlBn6PDFcfsDDWIHFs3 -x pBE2KvVaYN_569FAGOeEKr,8PXwCwkVb7L_2huJ6o8FBgE -y BVBszEaZArdp0XuVWgxQkDUSr6ylDCgi5t2i7nH4nMEh9ehZeyrPSbs8OnOu_pENXjVJguH45QnxV4k6P5gkxfhm989b23uMqVCd5dWc,uvzHiQfto4zzEJvjhKhCd30XIosCKe2XR3Uh8ULKXULqgA9skmQqlQ0yTWHk_vf0gf0Elij0iedGYRKHrbpTg5CnYvItrt1J9NalREqy
this is a secret.同时可以使用中文。
````

或者：
````
lhx@DESKTOP-0GALLEM:~$ shamir decrypt -n 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_MV0yI9PnICUSbyxh9iceG0W0HlBn6PDFcfsDDWIHFs3 -x pBE2KvVaYN_569FAGOeEKr -y BVBszEaZArdp0XuVWgxQkDUSr6ylDCgi5t2i7nH4nMEh9ehZeyrPSbs8OnOu_pENXjVJguH45QnxV4k6P5gkxfhm989b23uMqVCd5dWc,uvzHiQfto4zzEJvjhKhCd30XIosCKe2XR3Uh8ULKXULqgA9skmQqlQ0yTWHk_vf0gf0Elij0iedGYRKHrbpTg5CnYvItrt1J9NalREqy  -x 8PXwCwkVb7L_2huJ6o8FBgE
this is a secret.同时可以使用中文。
````

# 详细介绍：
## 背景
现在大家很熟悉的加解密场景，可以看做A hash成 B，由B可以还原成原来的A

但是在某些场景下，需要将密钥拆分成许多子密钥，让每一方分别持有一个子密钥；需要密钥时，必须有一定数量的子密钥合在一起才能恢复出整个密钥。这样的过程称为密钥共享，更广义的，也可称为秘密共享

使用秘密共享方案进行密钥管理，有如下优点：
1. 解决主密钥丢失即所有密钥丢失的问题
2. 避免出现权限最高的个体
3. 提高了密钥的安全性，因为攻击者必须获得足够多的子密钥才能恢复出主密钥

所有秘密共享方案中，最经典的要数Shamir于1979年提出的 **(k,n)门限秘密共享方案**

## 方案简述：
- 将秘密S分成n个子秘密
- 任何 k 个子秘密合在一起能恢复出秘密S
- 任何 <k 个子秘密合在一起无法恢复出秘密S

## 方案细节：
### 加密：
公共参数：参与方数n, 门限值k ≤ n, 素数p > n

秘密S：有限域Fp = {0,1,…,p-1}中的数

秘密分发：随机选取Fp上的k-1次多项式

![image](https://user-images.githubusercontent.com/40929503/198340494-02f984b8-6003-42c1-abbe-02f8ca9babb8.png)

再令a0 = S

任取n个Fp中不同的非零的数x1,…,xn,分别求出y1 = f(x1), …, yn = f(xn)

n份子秘密分别为(x1, y1) , … , (xn, yn)

### 恢复秘密：
任意给定k个点, 不妨就设为前k个点x1,…,xk, 首先构造k个基多项式g1,…,gk,

![image](https://user-images.githubusercontent.com/40929503/198340814-0f89bf48-bd0b-4978-a77a-bd2db2471ea6.png)

gi 满足gi(xi) = 1, 且gi(xj) = 0( j ≠ i )

计算多项式![image](https://user-images.githubusercontent.com/40929503/198340876-79ebe663-4e26-4357-bc80-562392f01e7c.png), 则f'(x)在xi点的取值就是yi  

而多项式f(x)在xi点的取值也是yi 

f'(x)与f(x)都不超过k-1次, 又在k个点取值一样，所以f(x) = f'(x)   

最后计算f'(0) = f(0) = a0 = S，恢复秘密S


