## Performance issues

### Problem

Only indexing in a for loop without any goroutines is extremely slow.
Creating a goroutine for each file does not work due to the fact that my system already has 25k man files and go limits go routines to 10k.

Without creating a preview via pandoc:

```
manbib master M $ go run . i
2023/05/16 14:57:43 starting index creation, this may take a while
2023/05/16 14:57:43 found 24362 man pages, starting indexing
2023/05/16 15:04:06 processed file: [24361/24362] 100.00 %
2023/05/16 15:04:06 done, took:  6m22.497429392s
```

To fix this i have to somehow split up the workload onto all the available go routines, because 6minutes is unacceptable.

With a fewer pages and creating a preview with pandoc:

```
2023/05/17 08:32:49 starting index creation, this may take a while
2023/05/17 08:32:49 found 904 man pages, starting indexing
2023/05/17 08:35:45 processed file: [903/904] 99.89 %
2023/05/17 08:35:45 done, took:  2m56.304404478s
```

Assume input `n = 24361` and available routines `g = 9000`.
By simply dividing `n` by `g`, we can calculate how many files every gorouting should process - `n / g = nPg`, where `nPg` means `n` per `g`.

### Python poc:

```python
import math

n = 24362
t = 9000

li = 0
i = 0

nPt = math.ceil( n / t)
while i < n:
    print(f"{li}->{i}")
    li = i
    i += nPt
print("nPt",nPt)
```
