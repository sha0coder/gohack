
# Build

    make

# Examples

pentest specific parameter
    ./carnivore -url 'https://testurl.com/index.php?id=2&mode=test' -go 3 -p mode


pentest all the parameters
    ./carnivore -url 'https://testurl.com/index.php?id=2&mode=test' -go 3 


pentest all parameters and also identify new hidden ones (slow)
    ./carnivore -url 'https://testurl.com/index.php?id=2&mode=test' -go 3 -new



# Pentest 

- test new parameters not found in url
- allow tests only spcecific parameter or all params
- many strategic payloads
- expert system recognize errors
- gauss to false positive reduction