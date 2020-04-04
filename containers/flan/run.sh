#!/bin/sh

current_time=$(date "+%Y.%m.%d")
root_dir=/shared/

mkdir -pv "${root_dir}xml_files"
mkdir -pv "${root_dir}reports"
mkdir -pv "${root_dir}$xml_dir"

xml_dir=xml_files/$current_time
report_file=reports/report_$current_time

mkdir -pv $root_dir$xml_dir

function upload {
    if [[ -z $upload ]]
    then
        return
    elif [ $upload = "aws" ]
    then
        python /aws_push.py $1
    elif [ $upload = "gcp" ]
    then
        python /gcp_push.py $1
    fi
}

function get_filename(){
    echo $1 | tr / -
}

while IFS= read -r line; do
  [ -z "$line" ] && continue

  filename=$(get_filename $line)
  filepath="$root_dir$xml_dir/$filename"

  if test -f "${filepath}.xml"; then
    echo "Host already scanned today ($current_time), skipping scan"
  else
    echo "Scanning host: $line"
    nmap -sV -oX "${filepath}.xml" -oN - -v1 $@ --script=vulners/vulners.nse $line > "${filepath}.txt" 2>&1
  fi
  upload $xml_dir
done < /shared/ips.txt

python /output_report.py --input $root_dir$xml_dir --output $root_dir$report_file --ipfile /shared/ips.txt
sed -i 's/_/\\_/g' $root_dir$report_file.tex
sed -i 's/\$/\\\$/g' $root_dir$report_file.tex
sed -i 's/#/\\#/g' $root_dir$report_file.tex
sed -i 's/%/\\%/g' $root_dir$report_file.tex
upload $report_file.tex
