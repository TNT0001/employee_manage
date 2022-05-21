#!/bin/bash

DW_INTERNAL_PATH='internal'

[ -d /migrations ] || mkdir /migrations

for d in "$DW_INTERNAL_PATH"/modules/* ; do
	filename=$(basename "$d")
	if [ "$filename" = "*" ]; then
	    continue
	fi
	mv "$DW_INTERNAL_PATH"/modules/"$filename"/migrations /migrations/"$filename";
		if [ "$?" -ne 0 ]; then
            echo "fail to mv migrations file $DW_INTERNAL_PATH/modules/$filename/migrations"
            exit 1
        fi;
done

[ -d /schemas ] || mkdir /schemas

for d in "$DW_INTERNAL_PATH"/service/* ; do
	filename=$(basename "$d")
	if [ "$filename" = "*" ]; then
	    continue
	fi
	mv "$DW_INTERNAL_PATH"/service/"$filename"/schemas /schemas/"$filename"
	if [ "$?" -ne 0 ]; then
        echo "fail to mv validate file $DW_INTERNAL_PATH/service/$filename/schemas"
        exit 1
    fi;
done