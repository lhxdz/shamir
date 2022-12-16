# shamir-tools
this is a tool for shamir $(k, n)$ threshold scheme

# 使用方式
快速安装：

````
git clone https://github.com/lhxdz/shamir.git
cd shamir
make install
````

使用 `shamir --version` 来验证

## 加密：支持加密命令行输入字符串、加密指定文件
输入门限值 $t$ 、密钥个数 $n$ 、秘密 即可加密，将会生成一个必须密钥 $necessary\\_key$ 、 $n$ 个密钥(每个密钥包含 $x$ 、 $y$ )，其中任意 $t$ 个密钥可以恢复秘密

````
root@DESKTOP-0GALLEM:~/project/shamir-tools# shamir encrypt -t 2 -n 3 "this is a secret.同时可以使用中文。"
necessary key: 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_41HtWR0abOmnN2ZbElB9ojNixNHGK2ZVIXPpOHs8nffiT
+-------------------------+-------------------------------------------------------------------------------------------------------------+
|          KEY X          |                                                    KEY Y                                                    |
+-------------------------+-------------------------------------------------------------------------------------------------------------+
| 26NJWXnvHHD_5NG3WEZJY6c | 1qXaP2vhGW7j7hDFBvdwP5fyMxFUJFZKCNkPHJg5e1XhkM4X6dZaKlGRd4X2F_13CXA5E6TcTMJGZumnJdRrnokzw97X65T2M8HPmHN3oQ3 |
+-------------------------+-------------------------------------------------------------------------------------------------------------+
| 4jXIgA9Eugh_3RsPU0ttqvx | 4DkyYRocvXq46aZbQt79AHoIvRaSJ5gUuHlv6uh6EvwIn01mtf5GjPdfiRvcP_3LW8yin7eBVHjw7Th5Lo0f9QaUBUsyeLRvUT6sflsmfiq |
+-------------------------+-------------------------------------------------------------------------------------------------------------+
| 8aLeam8AWrI_3WJ6UTFhDUl | 1aGTxpXJMgfEs6x7scFeK3Al5H1H0Dm13miSdJESP78YcptNVm7G2KmW7tRun_XFOBNvqZnBI43NP5uFw3rnHDpifKQ6ebCd1m5fWTBFEA  |
+-------------------------+-------------------------------------------------------------------------------------------------------------+
````

## 解密：支持从命令行获取密钥解密、从指定文件夹读取密钥文件解密
分别输入 $t$ 个密钥 $x$ 和 $y$ ，其中第 $i$ 个 $x$ 和第 $i$ 个 $y$ 是一对密钥，同时输入 $necessary\\_key$ ，即可恢复秘密

````
lhx@DESKTOP-0GALLEM:~$ shamir decrypt -n 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_41HtWR0abOmnN2ZbElB9ojNixNHGK2ZVIXPpOHs8nffiT -x 26NJWXnvHHD_5NG3WEZJY6c,4jXIgA9Eugh_3RsPU0ttqvx -y 1qXaP2vhGW7j7hDFBvdwP5fyMxFUJFZKCNkPHJg5e1XhkM4X6dZaKlGRd4X2F_13CXA5E6TcTMJGZumnJdRrnokzw97X65T2M8HPmHN3oQ3,4DkyYRocvXq46aZbQt79AHoIvRaSJ5gUuHlv6uh6EvwIn01mtf5GjPdfiRvcP_3LW8yin7eBVHjw7Th5Lo0f9QaUBUsyeLRvUT6sflsmfiq
this is a secret.同时可以使用中文。
````

或者：
````
lhx@DESKTOP-0GALLEM:~$ shamir decrypt -n 6IznUBFJvXlEgv1jFCGH4OsE4zPmhPVcHvcyrzVXPOMjCuu1gZxZgbq1CAWmd_41HtWR0abOmnN2ZbElB9ojNixNHGK2ZVIXPpOHs8nffiT -x 4jXIgA9Eugh_3RsPU0ttqvx -y 4DkyYRocvXq46aZbQt79AHoIvRaSJ5gUuHlv6uh6EvwIn01mtf5GjPdfiRvcP_3LW8yin7eBVHjw7Th5Lo0f9QaUBUsyeLRvUT6sflsmfiq,1aGTxpXJMgfEs6x7scFeK3Al5H1H0Dm13miSdJESP78YcptNVm7G2KmW7tRun_XFOBNvqZnBI43NP5uFw3rnHDpifKQ6ebCd1m5fWTBFEA  -x 8aLeam8AWrI_3WJ6UTFhDUl
this is a secret.同时可以使用中文。
````

**更多使用方式，请使用 `shamir --help`**

# 详细介绍：
## 背景
现在大家很熟悉的加解密场景，可以看做A hash成 B，由B可以还原成原来的A

但是在某些场景下，需要将密钥拆分成许多子密钥，让每一方分别持有一个子密钥；需要密钥时，必须有一定数量的子密钥合在一起才能恢复出整个密钥。这样的过程称为密钥共享，更广义的，也可称为秘密共享

使用秘密共享方案进行密钥管理，有如下优点：
1. 解决主密钥丢失即所有密钥丢失的问题
2. 避免出现权限最高的个体
3. 提高了密钥的安全性，因为攻击者必须获得足够多的子密钥才能恢复出主密钥

所有秘密共享方案中，最经典的要数Shamir于1979年提出的 **$(k,n)$ 门限秘密共享方案**

## 方案简述：
- 将秘密 $S$ 分成 $n$ 个子秘密
- 任何 $k$ 个子秘密合在一起能恢复出秘密 $S$
- 任何 $\textless k$ 个子秘密合在一起无法恢复出秘密 $S$

## 方案细节：
### 加密：
公共参数：参与方数 $n$ ，门限值 $k \leq n$ , 素数 $p>n$
秘密 $S$ ：有限域 $Fp=\{0,1,\dots,p-1\}$ 中的数

秘密分发：随机选取 $Fp$ 上的 $k-1$ 次多项式

$$
f(x)=(a_0+a_1+\cdots+a_{k-1}x^{k-1})\ mod\ p
$$

再令 $a_0=S$

任取 $n$ 个 $Fp$ 中不同的非零的数 $x_1,\dots,x_n$ ,分别求出 $y_1=f(x_1),\dots,y_n=f(x_n)$

$n$ 份子秘密分别为 $(x_1,y_1),\dots,(x_n,y_n)$

### 恢复秘密：
任意给定 $k$ 个点，不妨就设为前 $k$ 个点 $k_1,\dots,x_k$ ，首先构造 $k$ 个基多项式 $g_1,\dots,g_k$ ，

$$
g_i(x)=\prod^k_{j=1,j\neq i}\frac{(x-x_j)}{(x_i-x_j)}\ mod\ p
$$

$g_i$ 满足 $g_i(x_i)=1$ ，且 $g_i(x_j)=0(j\neq i)$

计算多项式

$$
f'(x)=\sum^k_{i=1}y_i\cdot g_i(x)\ mod\ p
$$

则 $f'(x)$ 在 $x_i$ 点的取值就是 $y_i$


而多项式 $f(x)$ 在 $x_i$ 点的取值也是 $y_i$

$f'(x)$ 与 $f(x)$ 都不超过 $k-1$ 次, 又在 $k$ 个点取值一样，所以 $f(x) = f'(x)$

最后计算 $f'(0) = f(0) = a_0 = S$ ，恢复秘密 $S$
