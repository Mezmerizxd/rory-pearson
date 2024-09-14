# Check if Go works
if ! command -v go &> /dev/null
then
    echo "go could not be found"
    exit
else
    echo "[)(] Found Go"
fi

# Check if yarn works
if ! command -v yarn &> /dev/null
then
    echo "yarn could not be found"
    exit
else
    echo "[)(] Found Yarn"
fi

# Install UI dependencies in the ui folder
cd ui
echo "[)(] Installing UI dependencies"
yarn install

## Build the UI
echo "[)(] Building UI"
yarn build

# Go back to the root folder
cd ..