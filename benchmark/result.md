
### 基准
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	  137570	     42951 ns/op	    4088 B/op	       9 allocs/op
BenchmarkBinarySearch-4           	228617994	        26.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkBubbleSort-4             	   44918	    133785 ns/op	    4088 B/op	       9 allocs/op
BenchmarkShuffle-4                	  330484	     18054 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringCompress-4         	    3034	   1903760 ns/op	  876214 B/op	      33 allocs/op
BenchmarkEncryptAndDecrypt-4      	  590178	      9990 ns/op	    1312 B/op	      10 allocs/op
BenchmarkWrite2DevNull-4          	 1428777	      4202 ns/op	     304 B/op	       5 allocs/op
BenchmarkWrite2TmpFile-4          	  535009	     10967 ns/op	     128 B/op	       1 allocs/op
BenchmarkLocalPipeReadWrite-4     	  265272	     21792 ns/op	    2176 B/op	      18 allocs/op
BenchmarkSerialNetRequest-4       	     387	  15407760 ns/op	   40489 B/op	     480 allocs/op
BenchmarkConcurrentNetRequest-4   	    1713	   3576828 ns/op	  136009 B/op	     990 allocs/op
PASS
ok  	benchmark	76.928s
```
### 采样率 0

PASS
ok  	github.com/lvyahui8/ellyn/instr	0.227s
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	  106446	     54818 ns/op	    4089 B/op	       9 allocs/op
BenchmarkBinarySearch-4           	82452814	        72.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkBubbleSort-4             	   32930	    181709 ns/op	    4092 B/op	       9 allocs/op
BenchmarkShuffle-4                	  331177	     18059 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringCompress-4         	    2541	   2289385 ns/op	  876263 B/op	      35 allocs/op
BenchmarkEncryptAndDecrypt-4      	  523080	     10913 ns/op	    1344 B/op	      12 allocs/op
BenchmarkWrite2DevNull-4          	 1373026	      4384 ns/op	     304 B/op	       5 allocs/op
BenchmarkWrite2TmpFile-4          	  527662	     11293 ns/op	     128 B/op	       1 allocs/op
BenchmarkLocalPipeReadWrite-4     	  270783	     20607 ns/op	    2193 B/op	      18 allocs/op
BenchmarkSerialNetRequest-4       	     434	  13748094 ns/op	   40883 B/op	     494 allocs/op
BenchmarkConcurrentNetRequest-4   	    1784	   3561672 ns/op	  136328 B/op	     992 allocs/op
PASS
ok  	benchmark	74.858s
```
性能差异
```text
goos: linux
goarch: arm64
pkg: benchmark
                       │   old.out    │                new0.out                │
                       │    sec/op    │    sec/op      vs base                 │
QuickSort-4              42.95µ ± ∞ ¹    54.82µ ± ∞ ¹        ~ (p=1.000 n=1) ²
BinarySearch-4           26.24n ± ∞ ¹    72.02n ± ∞ ¹        ~ (p=1.000 n=1) ²
BubbleSort-4             133.8µ ± ∞ ¹    181.7µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Shuffle-4                18.05µ ± ∞ ¹    18.06µ ± ∞ ¹        ~ (p=1.000 n=1) ²
StringCompress-4         1.904m ± ∞ ¹    2.289m ± ∞ ¹        ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      9.990µ ± ∞ ¹   10.913µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2DevNull-4          4.202µ ± ∞ ¹    4.384µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2TmpFile-4          10.97µ ± ∞ ¹    11.29µ ± ∞ ¹        ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     21.79µ ± ∞ ¹    20.61µ ± ∞ ¹        ~ (p=1.000 n=1) ²
SerialNetRequest-4       15.41m ± ∞ ¹    13.75m ± ∞ ¹        ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   3.577m ± ∞ ¹    3.562m ± ∞ ¹        ~ (p=1.000 n=1) ²
geomean                  47.63µ          55.75µ        +17.04%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                       │    old.out    │               new0.out                │
                       │     B/op      │     B/op       vs base                │
QuickSort-4              3.992Ki ± ∞ ¹   3.993Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4             0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
BubbleSort-4             3.992Ki ± ∞ ¹   3.996Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                  0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
StringCompress-4         855.7Ki ± ∞ ¹   855.7Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      1.281Ki ± ∞ ¹   1.312Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2DevNull-4            304.0 ± ∞ ¹     304.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2TmpFile-4            128.0 ± ∞ ¹     128.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
LocalPipeReadWrite-4     2.125Ki ± ∞ ¹   2.142Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       39.54Ki ± ∞ ¹   39.92Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   132.8Ki ± ∞ ¹   133.1Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
geomean                              ⁴                  +0.41%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                       │   old.out   │              new0.out               │
                       │  allocs/op  │  allocs/op   vs base                │
QuickSort-4              9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4           0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BubbleSort-4             9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
StringCompress-4         33.00 ± ∞ ¹   35.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
EncryptAndDecrypt-4      10.00 ± ∞ ¹   12.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2DevNull-4          5.000 ± ∞ ¹   5.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2TmpFile-4          1.000 ± ∞ ¹   1.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     18.00 ± ∞ ¹   18.00 ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       480.0 ± ∞ ¹   494.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
ConcurrentNetRequest-4   990.0 ± ∞ ¹   992.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
geomean                            ⁴                +2.50%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² all samples are equal
³ need >= 4 samples to detect a difference at alpha level 0.05
⁴ summaries must be >0 to compute geomean
```
### 采样率 0.0001

PASS
ok  	github.com/lvyahui8/ellyn/instr	0.246s
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	  107018	     55802 ns/op	    4089 B/op	       9 allocs/op
BenchmarkBinarySearch-4           	81365398	        72.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkBubbleSort-4             	   33294	    182848 ns/op	    4093 B/op	       9 allocs/op
BenchmarkShuffle-4                	  320906	     18466 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringCompress-4         	    2468	   3261636 ns/op	  876280 B/op	      35 allocs/op
BenchmarkEncryptAndDecrypt-4      	  563416	     10802 ns/op	    1344 B/op	      12 allocs/op
BenchmarkWrite2DevNull-4          	 1368524	      4353 ns/op	     304 B/op	       5 allocs/op
BenchmarkWrite2TmpFile-4          	  521224	     11328 ns/op	     128 B/op	       1 allocs/op
BenchmarkLocalPipeReadWrite-4     	  272166	     20679 ns/op	    2193 B/op	      18 allocs/op
BenchmarkSerialNetRequest-4       	     435	  13852948 ns/op	   40875 B/op	     494 allocs/op
BenchmarkConcurrentNetRequest-4   	    1730	   3471552 ns/op	  136226 B/op	     992 allocs/op
PASS
ok  	benchmark	77.277s
```
性能差异
```text
goos: linux
goarch: arm64
pkg: benchmark
                       │   old.out    │                new1.out                │
                       │    sec/op    │    sec/op      vs base                 │
QuickSort-4              42.95µ ± ∞ ¹    55.80µ ± ∞ ¹        ~ (p=1.000 n=1) ²
BinarySearch-4           26.24n ± ∞ ¹    72.32n ± ∞ ¹        ~ (p=1.000 n=1) ²
BubbleSort-4             133.8µ ± ∞ ¹    182.8µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Shuffle-4                18.05µ ± ∞ ¹    18.47µ ± ∞ ¹        ~ (p=1.000 n=1) ²
StringCompress-4         1.904m ± ∞ ¹    3.262m ± ∞ ¹        ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      9.990µ ± ∞ ¹   10.802µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2DevNull-4          4.202µ ± ∞ ¹    4.353µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2TmpFile-4          10.97µ ± ∞ ¹    11.33µ ± ∞ ¹        ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     21.79µ ± ∞ ¹    20.68µ ± ∞ ¹        ~ (p=1.000 n=1) ²
SerialNetRequest-4       15.41m ± ∞ ¹    13.85m ± ∞ ¹        ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   3.577m ± ∞ ¹    3.472m ± ∞ ¹        ~ (p=1.000 n=1) ²
geomean                  47.63µ          57.69µ        +21.10%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                       │    old.out    │               new1.out                │
                       │     B/op      │     B/op       vs base                │
QuickSort-4              3.992Ki ± ∞ ¹   3.993Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4             0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
BubbleSort-4             3.992Ki ± ∞ ¹   3.997Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                  0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
StringCompress-4         855.7Ki ± ∞ ¹   855.7Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      1.281Ki ± ∞ ¹   1.312Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2DevNull-4            304.0 ± ∞ ¹     304.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2TmpFile-4            128.0 ± ∞ ¹     128.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
LocalPipeReadWrite-4     2.125Ki ± ∞ ¹   2.142Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       39.54Ki ± ∞ ¹   39.92Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   132.8Ki ± ∞ ¹   133.0Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
geomean                              ⁴                  +0.41%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                       │   old.out   │              new1.out               │
                       │  allocs/op  │  allocs/op   vs base                │
QuickSort-4              9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4           0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BubbleSort-4             9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
StringCompress-4         33.00 ± ∞ ¹   35.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
EncryptAndDecrypt-4      10.00 ± ∞ ¹   12.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2DevNull-4          5.000 ± ∞ ¹   5.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2TmpFile-4          1.000 ± ∞ ¹   1.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     18.00 ± ∞ ¹   18.00 ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       480.0 ± ∞ ¹   494.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
ConcurrentNetRequest-4   990.0 ± ∞ ¹   992.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
geomean                            ⁴                +2.50%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² all samples are equal
³ need >= 4 samples to detect a difference at alpha level 0.05
⁴ summaries must be >0 to compute geomean
```
### 采样率 0.001

PASS
ok  	github.com/lvyahui8/ellyn/instr	0.260s
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	  109830	     55985 ns/op	    4093 B/op	       9 allocs/op
BenchmarkBinarySearch-4           	79786590	        73.72 ns/op	       0 B/op	       0 allocs/op
BenchmarkBubbleSort-4             	   32172	    183116 ns/op	    4095 B/op	       9 allocs/op
BenchmarkShuffle-4                	  326594	     18158 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringCompress-4         	    1932	   3094113 ns/op	  876283 B/op	      35 allocs/op
BenchmarkEncryptAndDecrypt-4      	  509672	     10636 ns/op	    1345 B/op	      12 allocs/op
BenchmarkWrite2DevNull-4          	 1378228	      4326 ns/op	     304 B/op	       5 allocs/op
BenchmarkWrite2TmpFile-4          	  525272	     11313 ns/op	     128 B/op	       1 allocs/op
BenchmarkLocalPipeReadWrite-4     	  288429	     20261 ns/op	    2195 B/op	      18 allocs/op
BenchmarkSerialNetRequest-4       	     430	  13803820 ns/op	   40914 B/op	     494 allocs/op
BenchmarkConcurrentNetRequest-4   	    1750	   3474480 ns/op	  136277 B/op	     992 allocs/op
PASS
ok  	benchmark	74.842s
```
性能差异
```text
goos: linux
goarch: arm64
pkg: benchmark
                       │   old.out    │                new2.out                │
                       │    sec/op    │    sec/op      vs base                 │
QuickSort-4              42.95µ ± ∞ ¹    55.98µ ± ∞ ¹        ~ (p=1.000 n=1) ²
BinarySearch-4           26.24n ± ∞ ¹    73.72n ± ∞ ¹        ~ (p=1.000 n=1) ²
BubbleSort-4             133.8µ ± ∞ ¹    183.1µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Shuffle-4                18.05µ ± ∞ ¹    18.16µ ± ∞ ¹        ~ (p=1.000 n=1) ²
StringCompress-4         1.904m ± ∞ ¹    3.094m ± ∞ ¹        ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      9.990µ ± ∞ ¹   10.636µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2DevNull-4          4.202µ ± ∞ ¹    4.326µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2TmpFile-4          10.97µ ± ∞ ¹    11.31µ ± ∞ ¹        ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     21.79µ ± ∞ ¹    20.26µ ± ∞ ¹        ~ (p=1.000 n=1) ²
SerialNetRequest-4       15.41m ± ∞ ¹    13.80m ± ∞ ¹        ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   3.577m ± ∞ ¹    3.474m ± ∞ ¹        ~ (p=1.000 n=1) ²
geomean                  47.63µ          57.21µ        +20.10%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                       │    old.out    │               new2.out                │
                       │     B/op      │     B/op       vs base                │
QuickSort-4              3.992Ki ± ∞ ¹   3.997Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4             0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
BubbleSort-4             3.992Ki ± ∞ ¹   3.999Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                  0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
StringCompress-4         855.7Ki ± ∞ ¹   855.7Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      1.281Ki ± ∞ ¹   1.313Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2DevNull-4            304.0 ± ∞ ¹     304.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2TmpFile-4            128.0 ± ∞ ¹     128.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
LocalPipeReadWrite-4     2.125Ki ± ∞ ¹   2.144Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       39.54Ki ± ∞ ¹   39.96Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   132.8Ki ± ∞ ¹   133.1Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
geomean                              ⁴                  +0.45%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                       │   old.out   │              new2.out               │
                       │  allocs/op  │  allocs/op   vs base                │
QuickSort-4              9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4           0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BubbleSort-4             9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
StringCompress-4         33.00 ± ∞ ¹   35.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
EncryptAndDecrypt-4      10.00 ± ∞ ¹   12.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2DevNull-4          5.000 ± ∞ ¹   5.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2TmpFile-4          1.000 ± ∞ ¹   1.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     18.00 ± ∞ ¹   18.00 ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       480.0 ± ∞ ¹   494.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
ConcurrentNetRequest-4   990.0 ± ∞ ¹   992.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
geomean                            ⁴                +2.50%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² all samples are equal
³ need >= 4 samples to detect a difference at alpha level 0.05
⁴ summaries must be >0 to compute geomean
```
### 采样率 0.01

PASS
ok  	github.com/lvyahui8/ellyn/instr	0.288s
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	  107619	     56872 ns/op	    4097 B/op	       9 allocs/op
BenchmarkBinarySearch-4           	68519330	        87.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkBubbleSort-4             	   32450	    186070 ns/op	    4099 B/op	       9 allocs/op
BenchmarkShuffle-4                	  329918	     18133 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringCompress-4         	    2487	   2748280 ns/op	  876314 B/op	      35 allocs/op
BenchmarkEncryptAndDecrypt-4      	  545625	     11146 ns/op	    1347 B/op	      12 allocs/op
BenchmarkWrite2DevNull-4          	 1302597	      4551 ns/op	     304 B/op	       5 allocs/op
BenchmarkWrite2TmpFile-4          	  510146	     11407 ns/op	     128 B/op	       1 allocs/op
BenchmarkLocalPipeReadWrite-4     	  284913	     21222 ns/op	    2198 B/op	      18 allocs/op
BenchmarkSerialNetRequest-4       	     433	  13815425 ns/op	   40897 B/op	     494 allocs/op
BenchmarkConcurrentNetRequest-4   	    1791	   3451836 ns/op	  136431 B/op	     993 allocs/op
PASS
ok  	benchmark	77.021s
```
性能差异
```text
goos: linux
goarch: arm64
pkg: benchmark
                       │   old.out    │                new3.out                │
                       │    sec/op    │    sec/op      vs base                 │
QuickSort-4              42.95µ ± ∞ ¹    56.87µ ± ∞ ¹        ~ (p=1.000 n=1) ²
BinarySearch-4           26.24n ± ∞ ¹    87.46n ± ∞ ¹        ~ (p=1.000 n=1) ²
BubbleSort-4             133.8µ ± ∞ ¹    186.1µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Shuffle-4                18.05µ ± ∞ ¹    18.13µ ± ∞ ¹        ~ (p=1.000 n=1) ²
StringCompress-4         1.904m ± ∞ ¹    2.748m ± ∞ ¹        ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      9.990µ ± ∞ ¹   11.146µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2DevNull-4          4.202µ ± ∞ ¹    4.551µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2TmpFile-4          10.97µ ± ∞ ¹    11.41µ ± ∞ ¹        ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     21.79µ ± ∞ ¹    21.22µ ± ∞ ¹        ~ (p=1.000 n=1) ²
SerialNetRequest-4       15.41m ± ∞ ¹    13.82m ± ∞ ¹        ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   3.577m ± ∞ ¹    3.452m ± ∞ ¹        ~ (p=1.000 n=1) ²
geomean                  47.63µ          58.41µ        +22.62%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                       │    old.out    │               new3.out                │
                       │     B/op      │     B/op       vs base                │
QuickSort-4              3.992Ki ± ∞ ¹   4.001Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4             0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
BubbleSort-4             3.992Ki ± ∞ ¹   4.003Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                  0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
StringCompress-4         855.7Ki ± ∞ ¹   855.8Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      1.281Ki ± ∞ ¹   1.315Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2DevNull-4            304.0 ± ∞ ¹     304.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2TmpFile-4            128.0 ± ∞ ¹     128.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
LocalPipeReadWrite-4     2.125Ki ± ∞ ¹   2.146Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       39.54Ki ± ∞ ¹   39.94Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   132.8Ki ± ∞ ¹   133.2Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
geomean                              ⁴                  +0.50%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                       │   old.out   │              new3.out               │
                       │  allocs/op  │  allocs/op   vs base                │
QuickSort-4              9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4           0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BubbleSort-4             9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
StringCompress-4         33.00 ± ∞ ¹   35.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
EncryptAndDecrypt-4      10.00 ± ∞ ¹   12.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2DevNull-4          5.000 ± ∞ ¹   5.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2TmpFile-4          1.000 ± ∞ ¹   1.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     18.00 ± ∞ ¹   18.00 ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       480.0 ± ∞ ¹   494.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
ConcurrentNetRequest-4   990.0 ± ∞ ¹   993.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
geomean                            ⁴                +2.51%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² all samples are equal
³ need >= 4 samples to detect a difference at alpha level 0.05
⁴ summaries must be >0 to compute geomean
```
### 采样率 0.1

PASS
ok  	github.com/lvyahui8/ellyn/instr	0.225s
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	   92020	     66116 ns/op	    4100 B/op	       9 allocs/op
BenchmarkBinarySearch-4           	27196087	       224.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBubbleSort-4             	   26347	    228685 ns/op	    4102 B/op	       9 allocs/op
BenchmarkShuffle-4                	  319908	     18743 ns/op	       0 B/op	       0 allocs/op
BenchmarkStringCompress-4         	    2794	   2768643 ns/op	  876625 B/op	      38 allocs/op
BenchmarkEncryptAndDecrypt-4      	  515992	     11863 ns/op	    1353 B/op	      12 allocs/op
BenchmarkWrite2DevNull-4          	 1258263	      4682 ns/op	     306 B/op	       5 allocs/op
BenchmarkWrite2TmpFile-4          	  504908	     11712 ns/op	     129 B/op	       1 allocs/op
BenchmarkLocalPipeReadWrite-4     	  291871	     21474 ns/op	    2202 B/op	      18 allocs/op
BenchmarkSerialNetRequest-4       	     422	  13805941 ns/op	   40968 B/op	     494 allocs/op
BenchmarkConcurrentNetRequest-4   	    1776	   3452681 ns/op	  136531 B/op	     996 allocs/op
PASS
ok  	benchmark	78.926s
```
性能差异
```text
goos: linux
goarch: arm64
pkg: benchmark
                       │   old.out    │                new4.out                │
                       │    sec/op    │    sec/op      vs base                 │
QuickSort-4              42.95µ ± ∞ ¹    66.12µ ± ∞ ¹        ~ (p=1.000 n=1) ²
BinarySearch-4           26.24n ± ∞ ¹   224.10n ± ∞ ¹        ~ (p=1.000 n=1) ²
BubbleSort-4             133.8µ ± ∞ ¹    228.7µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Shuffle-4                18.05µ ± ∞ ¹    18.74µ ± ∞ ¹        ~ (p=1.000 n=1) ²
StringCompress-4         1.904m ± ∞ ¹    2.769m ± ∞ ¹        ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      9.990µ ± ∞ ¹   11.863µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2DevNull-4          4.202µ ± ∞ ¹    4.682µ ± ∞ ¹        ~ (p=1.000 n=1) ²
Write2TmpFile-4          10.97µ ± ∞ ¹    11.71µ ± ∞ ¹        ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     21.79µ ± ∞ ¹    21.47µ ± ∞ ¹        ~ (p=1.000 n=1) ²
SerialNetRequest-4       15.41m ± ∞ ¹    13.81m ± ∞ ¹        ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   3.577m ± ∞ ¹    3.453m ± ∞ ¹        ~ (p=1.000 n=1) ²
geomean                  47.63µ          66.74µ        +40.11%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                       │    old.out    │               new4.out                │
                       │     B/op      │     B/op       vs base                │
QuickSort-4              3.992Ki ± ∞ ¹   4.004Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4             0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
BubbleSort-4             3.992Ki ± ∞ ¹   4.006Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                  0.000 ± ∞ ¹     0.000 ± ∞ ¹       ~ (p=1.000 n=1) ³
StringCompress-4         855.7Ki ± ∞ ¹   856.1Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      1.281Ki ± ∞ ¹   1.321Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2DevNull-4            304.0 ± ∞ ¹     306.0 ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2TmpFile-4            128.0 ± ∞ ¹     129.0 ± ∞ ¹       ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     2.125Ki ± ∞ ¹   2.150Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       39.54Ki ± ∞ ¹   40.01Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   132.8Ki ± ∞ ¹   133.3Ki ± ∞ ¹       ~ (p=1.000 n=1) ²
geomean                              ⁴                  +0.72%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ all samples are equal
⁴ summaries must be >0 to compute geomean

                       │   old.out   │              new4.out               │
                       │  allocs/op  │  allocs/op   vs base                │
QuickSort-4              9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BinarySearch-4           0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
BubbleSort-4             9.000 ± ∞ ¹   9.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Shuffle-4                0.000 ± ∞ ¹   0.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
StringCompress-4         33.00 ± ∞ ¹   38.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
EncryptAndDecrypt-4      10.00 ± ∞ ¹   12.00 ± ∞ ¹       ~ (p=1.000 n=1) ³
Write2DevNull-4          5.000 ± ∞ ¹   5.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
Write2TmpFile-4          1.000 ± ∞ ¹   1.000 ± ∞ ¹       ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     18.00 ± ∞ ¹   18.00 ± ∞ ¹       ~ (p=1.000 n=1) ²
SerialNetRequest-4       480.0 ± ∞ ¹   494.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
ConcurrentNetRequest-4   990.0 ± ∞ ¹   996.0 ± ∞ ¹       ~ (p=1.000 n=1) ³
geomean                            ⁴                +3.31%               ⁴
¹ need >= 6 samples for confidence interval at level 0.95
² all samples are equal
³ need >= 4 samples to detect a difference at alpha level 0.05
⁴ summaries must be >0 to compute geomean
```
### 采样率 1

PASS
ok  	github.com/lvyahui8/ellyn/instr	0.262s
```text
goos: linux
goarch: arm64
pkg: benchmark
BenchmarkQuickSort-4              	   38869	    153541 ns/op	    4129 B/op	      12 allocs/op
BenchmarkBinarySearch-4           	 3629982	      1644 ns/op	       9 B/op	       1 allocs/op
BenchmarkBubbleSort-4             	    9901	    604325 ns/op	    4121 B/op	      11 allocs/op
BenchmarkShuffle-4                	  261866	     22736 ns/op	       3 B/op	       1 allocs/op
BenchmarkStringCompress-4         	    2941	   2751838 ns/op	  877206 B/op	      44 allocs/op
BenchmarkEncryptAndDecrypt-4      	  337108	     17787 ns/op	    1408 B/op	      16 allocs/op
BenchmarkWrite2DevNull-4          	  882126	      6945 ns/op	     332 B/op	       6 allocs/op
BenchmarkWrite2TmpFile-4          	  412818	     14474 ns/op	     133 B/op	       2 allocs/op
BenchmarkLocalPipeReadWrite-4     	  240080	     26798 ns/op	    2243 B/op	      21 allocs/op
BenchmarkSerialNetRequest-4       	     428	  13911929 ns/op	   41007 B/op	     506 allocs/op
BenchmarkConcurrentNetRequest-4   	    1640	   3576815 ns/op	  137025 B/op	    1027 allocs/op
PASS
ok  	benchmark	74.694s
```
性能差异
```text
goos: linux
goarch: arm64
pkg: benchmark
                       │   old.out    │                 new5.out                 │
                       │    sec/op    │     sec/op      vs base                  │
QuickSort-4              42.95µ ± ∞ ¹    153.54µ ± ∞ ¹         ~ (p=1.000 n=1) ²
BinarySearch-4           26.24n ± ∞ ¹   1644.00n ± ∞ ¹         ~ (p=1.000 n=1) ²
BubbleSort-4             133.8µ ± ∞ ¹     604.3µ ± ∞ ¹         ~ (p=1.000 n=1) ²
Shuffle-4                18.05µ ± ∞ ¹     22.74µ ± ∞ ¹         ~ (p=1.000 n=1) ²
StringCompress-4         1.904m ± ∞ ¹     2.752m ± ∞ ¹         ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      9.990µ ± ∞ ¹    17.787µ ± ∞ ¹         ~ (p=1.000 n=1) ²
Write2DevNull-4          4.202µ ± ∞ ¹     6.945µ ± ∞ ¹         ~ (p=1.000 n=1) ²
Write2TmpFile-4          10.97µ ± ∞ ¹     14.47µ ± ∞ ¹         ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     21.79µ ± ∞ ¹     26.80µ ± ∞ ¹         ~ (p=1.000 n=1) ²
SerialNetRequest-4       15.41m ± ∞ ¹     13.91m ± ∞ ¹         ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   3.577m ± ∞ ¹     3.577m ± ∞ ¹         ~ (p=1.000 n=1) ²
geomean                  47.63µ           107.8µ        +126.22%
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05

                       │    old.out    │             new5.out             │
                       │     B/op      │     B/op       vs base           │
QuickSort-4              3.992Ki ± ∞ ¹   4.032Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
BinarySearch-4             0.000 ± ∞ ¹     9.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
BubbleSort-4             3.992Ki ± ∞ ¹   4.024Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
Shuffle-4                  0.000 ± ∞ ¹     3.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
StringCompress-4         855.7Ki ± ∞ ¹   856.6Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      1.281Ki ± ∞ ¹   1.375Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
Write2DevNull-4            304.0 ± ∞ ¹     332.0 ± ∞ ¹  ~ (p=1.000 n=1) ²
Write2TmpFile-4            128.0 ± ∞ ¹     133.0 ± ∞ ¹  ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     2.125Ki ± ∞ ¹   2.190Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
SerialNetRequest-4       39.54Ki ± ∞ ¹   40.05Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   132.8Ki ± ∞ ¹   133.8Ki ± ∞ ¹  ~ (p=1.000 n=1) ²
geomean                              ³   1.648Ki        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean

                       │   old.out   │            new5.out             │
                       │  allocs/op  │  allocs/op    vs base           │
QuickSort-4              9.000 ± ∞ ¹   12.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
BinarySearch-4           0.000 ± ∞ ¹    1.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
BubbleSort-4             9.000 ± ∞ ¹   11.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
Shuffle-4                0.000 ± ∞ ¹    1.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
StringCompress-4         33.00 ± ∞ ¹    44.00 ± ∞ ¹  ~ (p=1.000 n=1) ²
EncryptAndDecrypt-4      10.00 ± ∞ ¹    16.00 ± ∞ ¹  ~ (p=1.000 n=1) ²
Write2DevNull-4          5.000 ± ∞ ¹    6.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
Write2TmpFile-4          1.000 ± ∞ ¹    2.000 ± ∞ ¹  ~ (p=1.000 n=1) ²
LocalPipeReadWrite-4     18.00 ± ∞ ¹    21.00 ± ∞ ¹  ~ (p=1.000 n=1) ²
SerialNetRequest-4       480.0 ± ∞ ¹    506.0 ± ∞ ¹  ~ (p=1.000 n=1) ²
ConcurrentNetRequest-4   990.0 ± ∞ ¹   1027.0 ± ∞ ¹  ~ (p=1.000 n=1) ²
geomean                            ³    15.47        ?
¹ need >= 6 samples for confidence interval at level 0.95
² need >= 4 samples to detect a difference at alpha level 0.05
³ summaries must be >0 to compute geomean
```
