#!/usr/bin/env bash

RESULT_CONFIG=/jupyter/config.py
JUPYTER_BIN=/usr/local/bin/jupyter

echo "Generating Jupyter configuration"

if [ -z "$JUPYTER_CONF_FILE" ]; then
    export JUPYTER_CONF_FILE=jupyter_base_config.py
fi
echo "Using $JUPYTER_CONF_FILE as base file"
cp $JUPYTER_CONF_FILE $RESULT_CONFIG

if [ -z "$JUPYTER_PORT" ]; then
    export JUPYTER_PORT=8888
fi
echo "c.NotebookApp.port = $JUPYTER_PORT" >> $RESULT_CONFIG

if [ -z "$JUPYTER_PASSWD" ]; then
    echo "c.NotebookApp.password_required = False" >> $RESULT_CONFIG
else
    JUPYTER_HASH=$(python3 -c "from notebook.auth import passwd; print(passwd('${JUPYTER_PASSWD}'))")
    echo "c.NotebookApp.password = $JUPYTER_HASH" >> $RESULT_CONFIG
fi

if [ -z "$JUPYTER_NBDIR" ]; then
    export JUPYTER_NBDIR=/jupyter/Notebooks
fi
echo "c.NotebookApp.notebook_dir = $JUPYTER_NBDIR" >> $RESULT_CONFIG

echo "Starting Jupyter notebook"
$JUPYTER_BIN notebook --config=$RESULT_CONFIG

