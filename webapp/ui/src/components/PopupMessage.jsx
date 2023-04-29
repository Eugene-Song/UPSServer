export const PopupM = ({ isVisible, message, onClose }) => {
    if (!isVisible) {
      return null;
    }
  
    return (
      <div className="modal">
        <div className="modal-content">
          <p>{message}</p>
          <button onClick={onClose}>Close</button>
        </div>
      </div>
    );
  };

  // export default PopupM;

  export const modalStyle = `
    .modal {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: rgba(0, 0, 0, 0.5);
      z-index: 1000;
    }
    .modal-content {
      background-color: #fff;
      padding: 2rem;
      border-radius: 4px;
      width: 80%;
      max-width: 500px;
    }
  `;
  