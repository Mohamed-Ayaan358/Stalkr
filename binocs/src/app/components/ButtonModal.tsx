import { useState } from "react";
import axios from "axios";

const ButtonModal = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    url: "",
    time: 0,
  });

  const openModal = () => {
    setIsModalOpen(true);
  };

  const closeModal = () => {
    setIsModalOpen(false);
  };

  const handleChange = (e: any) => {
    const value =
      e.target.name === "time" ? parseInt(e.target.value, 10) : e.target.value;

    setFormData({
      ...formData,
      [e.target.name]: value,
    });
  };
  const handleFormSubmit = async (e: any) => {
    e.preventDefault();

    try {
      // Instead of directly sending the request to the API server, send it to your Next.js API route
      const response = await axios.post("/api/add", formData);
      console.log(response.data);

      closeModal();
    } catch (error) {
      console.error("Error submitting form:", error);
    }
  };

  return (
    <>
      <button
        onClick={openModal}
        className="bg-blue-500 text-white px-4 py-2 rounded"
      >
        Add a website
      </button>

      {isModalOpen && (
        <div className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div
              className="fixed inset-0 transition-opacity"
              onClick={closeModal}
            >
              <div className="absolute inset-0 bg-gray-500 opacity-75"></div>
            </div>
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen"></span>
            &#8203;
            <div
              className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
              role="dialog"
              aria-modal="true"
              aria-labelledby="modal-headline"
            >
              <form onSubmit={handleFormSubmit}>
                <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                  <label
                    className="block text-gray-700 text-sm font-bold mb-2"
                    htmlFor="name"
                  >
                    Name:
                  </label>
                  <input
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    className="w-full border border-gray-300 p-2 rounded mb-3"
                  />
                  <label
                    className="block text-gray-700 text-sm font-bold mb-2"
                    htmlFor="url"
                  >
                    Url:
                  </label>
                  <input
                    type="text"
                    name="url"
                    value={formData.url}
                    onChange={handleChange}
                    className="w-full border border-gray-300 p-2 rounded mb-3"
                  />
                  <label
                    className="block text-gray-700 text-sm font-bold mb-2"
                    htmlFor="time"
                  >
                    Time
                  </label>
                  <input
                    type="number"
                    name="time"
                    value={formData.time}
                    onChange={handleChange}
                    className="w-full border border-gray-300 p-2 rounded mb-3"
                  />
                </div>
                <div className=" gap-5 bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                  <button
                    type="submit"
                    className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-500 text-white focus:outline-none focus:border-blue-700 focus:shadow-outline-blue active:bg-blue-800"
                  >
                    Submit
                  </button>
                  <button
                    onClick={closeModal}
                    type="button"
                    className=" w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-sm text-gray-700 hover:bg-gray-100 focus:outline-none focus:border-blue-300 focus:shadow-outline-blue active:bg-gray-200"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default ButtonModal;
