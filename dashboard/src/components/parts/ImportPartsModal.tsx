import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { downloadPartsTemplate } from '@/api/parts';

interface ImportPartsModalProps {
  open: boolean;
  onClose: () => void;
  onImport: (file: File) => void;
  isLoading: boolean;
}

export function ImportPartsModal({
  open,
  onClose,
  onImport,
  isLoading,
}: ImportPartsModalProps) {
  const [file, setFile] = React.useState<File | null>(null);
  const [dragActive, setDragActive] = React.useState(false);
  const fileInputRef = React.useRef<HTMLInputElement>(null);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const droppedFile = e.dataTransfer.files[0];
      if (droppedFile.type === 'text/csv' || droppedFile.name.endsWith('.csv')) {
        setFile(droppedFile);
      }
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
    }
  };

  const handleSubmit = () => {
    if (file) {
      onImport(file);
    }
  };

  const handleClose = () => {
    setFile(null);
    onClose();
  };

  return (
    <Modal open={open} onClose={handleClose} className="max-w-lg">
      <ModalHeader onClose={handleClose}>Import Parts</ModalHeader>
      <ModalBody>
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <p className="text-sm text-gray-600">
              Upload a CSV file with parts data.
            </p>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={downloadPartsTemplate}
            >
              Download Template
            </Button>
          </div>

          <div
            className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
              dragActive
                ? 'border-blue-500 bg-blue-50'
                : file
                ? 'border-green-500 bg-green-50'
                : 'border-gray-300 hover:border-gray-400'
            }`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            <input
              ref={fileInputRef}
              type="file"
              accept=".csv"
              onChange={handleFileChange}
              className="hidden"
            />

            {file ? (
              <div className="space-y-2">
                <div className="text-green-600">
                  <svg
                    className="w-12 h-12 mx-auto"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                </div>
                <p className="font-medium text-gray-900">{file.name}</p>
                <p className="text-sm text-gray-500">
                  {(file.size / 1024).toFixed(1)} KB
                </p>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setFile(null)}
                >
                  Remove
                </Button>
              </div>
            ) : (
              <div className="space-y-2">
                <div className="text-gray-400">
                  <svg
                    className="w-12 h-12 mx-auto"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                    />
                  </svg>
                </div>
                <p className="text-gray-600">
                  Drag and drop a CSV file here, or{' '}
                  <button
                    type="button"
                    className="text-blue-600 hover:underline"
                    onClick={() => fileInputRef.current?.click()}
                  >
                    browse
                  </button>
                </p>
                <p className="text-xs text-gray-500">CSV files only</p>
              </div>
            )}
          </div>

          <div className="bg-gray-50 rounded-lg p-4">
            <h4 className="text-sm font-medium text-gray-700 mb-2">
              CSV Format
            </h4>
            <p className="text-xs text-gray-600 mb-2">
              Required columns: <code className="bg-gray-200 px-1 rounded">sku</code>,{' '}
              <code className="bg-gray-200 px-1 rounded">name</code>
            </p>
            <p className="text-xs text-gray-600">
              Optional columns:{' '}
              <code className="bg-gray-200 px-1 rounded">category</code>,{' '}
              <code className="bg-gray-200 px-1 rounded">description</code>,{' '}
              <code className="bg-gray-200 px-1 rounded">unitCostCents</code>,{' '}
              <code className="bg-gray-200 px-1 rounded">supplier</code>,{' '}
              <code className="bg-gray-200 px-1 rounded">supplierSku</code>
            </p>
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button type="button" variant="outline" onClick={handleClose}>
          Cancel
        </Button>
        <Button
          type="button"
          disabled={!file || isLoading}
          onClick={handleSubmit}
        >
          {isLoading ? 'Importing...' : 'Import Parts'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
