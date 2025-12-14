function readfile(file,     tmp, save_rs)
{
    save_rs = RS
    RS = "^$"
    getline tmp < file
    close(file)
    RS = save_rs

    return tmp
}

# // Embed wasm_exec.js\
/^\s*<!-- EMBED [A-Za-z0-9._]+ -->\s*$/ {
    match($0, /<!-- EMBED ([A-Za-z0-9._]+) -->/, arr);
    vv = readfile(arr[1]);
    printf("%s", vv);
    next;
}

{ print; }