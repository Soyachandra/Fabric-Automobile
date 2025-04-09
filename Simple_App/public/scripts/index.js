document.addEventListener("DOMContentLoaded", function () {
    document.getElementById("addForm").addEventListener("submit", addData);
    document.getElementById("readForm").addEventListener("submit", readData);
});

const addData = async (event) => {
    event.preventDefault();

    const formData = new FormData(document.getElementById("addForm"));

    try {
        const response = await fetch("/api/car", {
            method: "POST",
            body: formData, // Send as FormData (includes file)
        });

        const data = await response.json();
        alert("Car Created: " + JSON.stringify(data));
    } catch (err) {
        alert("Error creating car.");
        console.error(err);
    }
};

const readData = async (event) => {
    event.preventDefault();
    const carId = document.getElementById("carIdInput").value.trim();

    if (!carId) {
        alert("Please enter a valid Car ID.");
        return;
    }

    try {
        const response = await fetch(`/api/car/${carId}`);
        const responseData = await response.json();

        document.getElementById("queryResult").innerHTML = `<pre>${JSON.stringify(responseData, null, 2)}</pre>`;
    } catch (err) {
        alert("Error retrieving car details.");
        console.error(err);
    }
};
