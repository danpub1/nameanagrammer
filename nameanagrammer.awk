# gawk -M -f nameanagrammer.awk -v source=JohnSmith -v mode=FL >tmp
function hash(s,   ii, rv)
{
    rv = 1
    for (ii = 1; ii <= length(s); ii++)
       rv *= prime[substr(s,ii,1)]

    return rv
}

function remainder(source, first,    ii, fnd)
{
    for (ii = 1; ii <= length(first); ii++)
    {
        fnd = index(source, substr(first, ii, 1))
        if (fnd == 0)
            return ""
        source = (fnd == 1 ? "" : substr(source, 1, fnd-1)) (fnd == length(source) ? "" : substr(source, fnd + 1))
    }

    return source;
}

BEGIN {
    prime["A"] = 2;
    prime["B"] = 3;
    prime["C"] = 5;
    prime["D"] = 7;
    prime["E"] = 11;
    prime["F"] = 13;
    prime["G"] = 17;
    prime["H"] = 19;
    prime["I"] = 23;
    prime["J"] = 29;
    prime["K"] = 31;
    prime["L"] = 37;
    prime["M"] = 41;
    prime["N"] = 43;
    prime["O"] = 47;
    prime["P"] = 53;
    prime["Q"] = 59;
    prime["R"] = 61;
    prime["S"] = 67;
    prime["T"] = 71;
    prime["U"] = 73;
    prime["V"] = 79;
    prime["W"] = 83;
    prime["X"] = 89;
    prime["Y"] = 97;
    prime["Z"] = 101;

    while ((getline line < "Names_2010Census.csv") > 0)
    {
        split(line, arr, ",")
        code = hash(arr[1])
        last[code] = last[code] (length(last[code]) > 0 ? "/" : "") arr[1]
    }
    close("Names_2010Census.csv")

    while ((getline line < "allfirst.csv") > 0)
    {
        split(line, arr, ",")
        line = toupper(arr[1])
        code = hash(line)
        if (!(code in allfirst))
            first[line] = 0
        allfirst[code] = allfirst[code] (length(allfirst[code]) > 0 ? "/" : "") line
    }
    close("allfirst.csv")

    asorti(first)

    if (length(source) == 0)
        source = "JohnSmith"

    if (length(mode) == 0)
        mode = "FL"

    source = toupper(source);

    if (mode == "FML")
    {
        for (idx = 1; idx <= length(first); idx++)
        {
            afirst = first[idx]
            alast = remainder(source, afirst)
            if (length(alast) > 1 && idx < length(first))
            {
                for (idx2 = idx + 1; idx2 <=length(first); idx2++)
                {
                    bfirst = first[idx2]
                    if (bfirst > afirst)
                    {
                        newlast = remainder(alast, bfirst)
                        if (length(newlast) > 0)
                        {
                            code = hash(newlast)
                            if (code in last)
                            {
                                found[code] = found[code] (length(found[code]) > 0 ? " | " : "") "[" allfirst[hash(afirst)] " " allfirst[hash(bfirst)] "]"
                            }
                        }
                    }
                }
            }
        }
    }
    else if (mode == "FMIL")
    {
        for (idx = 1; idx <= length(first); idx++)
        {
            afirst = first[idx]
            alast = remainder(source, afirst)
            if (length(alast) > 1)
            {
                for (ii = 1; ii < length(alast); ii++)
                {
                    mi = substr(alast, ii, 1)
                    if (index(alast, mi) == ii)
                    {
                        newlast = remainder(alast, mi)
                        code = hash(newlast)
                        if (code in last)
                        {
                            found[code] = found[code] (length(found[code]) > 0 ? " | "  : "") "[" allfirst[hash(afirst)] " " mi "]"
                        }
                    }
                }
            }
        }
    }
    else # "FL"
    {
        for (idx = 1; idx <= length(first); idx++)
        {
            afirst = first[idx];
            alast = remainder(source, afirst)
            if (length(alast) > 0)
            {
                code = hash(alast)
                if (code in last)
                {
                    found[code] = allfirst[hash(afirst)]
                }
            }
        }
    }
    
    if (mode == "FL")
    {
        while ((getline line < "Names_2010Census.csv") > 0)
        {
            split(line, arr, ",")
            lastwt[arr[1]] = arr[4] * 1
        }
        close("Names_2010Census.csv")

        while ((getline line < "allfirst.csv") > 0)
        {
            split(line, arr, ",")
            firstwt[toupper(arr[1])] = arr[2] * 1
        }
        close("allfirst.csv")

        print "type NameEntry struct {\n    weight float32\n    low, high int64\n}\n"

        print "lastwt := map[string]NameEntry{"
        for (aa in lastwt)
        {
            ii = hash(aa)
            ll = ii % (0xFFFFFFFFFFFFFFFF + 1);
            hh = ii / (0xFFFFFFFFFFFFFFFF + 1);
            printf("   \"%s\": {%.2f, 0x%016x, 0x%016x},\n", aa, lastwt[aa], ll, hh)
            #ii = sprintf("%032x", ii)
            #ll = sprintf("%016x%016x", hh, ll)
            #if (ii != ll)
            #   printf("Error at %s\n", aa)
        }
        print "}\n"

        print "firstwt := map[string]NameEntry{"
        for (aa in firstwt)
        {
            ii = hash(aa)
            ll = ii % (0xFFFFFFFFFFFFFFFF + 1);
            hh = ii / (0xFFFFFFFFFFFFFFFF + 1);
            printf("   \"%s\": {%.8f, 0x%016x, 0x%016x},\n", aa, firstwt[aa], ll, hh)
        }
        print "}"

        split("", out)
        for (acode in found)
        {
            split(found[acode], fnarr, "/")
            split(last[acode], lnarr, "/")
            for (afn in fnarr)
            {
                for (aln in lnarr)
                {
                    s = sprintf("%0.25f,%s %s", firstwt[fnarr[afn]]/100000 * lastwt[lnarr[aln]], fnarr[afn], lnarr[aln])
                    out[length(out)+1] = s
                }   
            }   
        }

        asort(out)
        for (ii = length(out); ii > 0; ii--)
            print out[ii]
    }
    else
    {
        for (acode in found)
        {
            print found[acode] " ==> " last[acode]
        }
    }
}