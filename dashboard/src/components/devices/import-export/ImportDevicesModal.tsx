import * as React from 'react';
import {
  Upload,
  Download,
  AlertCircle,
  CheckCircle2,
  X,
  File,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { cn } from '@/lib/utils';
import type { ImportResult } from '@/types/device';
import { downloadImportTemplate } from '@/api/devices';

interface ImportDevicesModalProps {
  open: boolean;
  onClose: () => void;
  onImport: (file: File) => Promise<ImportResult>;
  isLoading: boolean;
}

export function ImportDevicesModal({
  open,
  onClose,
  onImport,
  isLoading,
}: ImportDevicesModalProps) {
  const [file, setFile] = React.useState<File | null>(null);
  const [importResult, setImportResult] = React.useState<ImportResult | null>(null);
  const [showErrors, setShowErrors] = React.useState(false);
  const [dragOver, setDragOver] = React.useState(false);
  const fileInputRef = React.useRef<HTMLInputElement>(null);

  // Reset state when modal opens
  React.useEffect(() => {
    if (open) {
      setFile(null);
      setImportResult(null);
      setShowErrors(false);
    }
  }, [open]);

  // Handle file drop
  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile && droppedFile.type === 'text/csv') {
      setFile(droppedFile);
      setImportResult(null);
    }
  };

  // Handle file select
  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setImportResult(null);
    }
  };

  // Handle import
  const handleImport = async () => {
    if (!file) return;

    const result = await onImport(file);
    setImportResult(result);
  };

  return (
    <Modal open={open} onClose={onClose} className="max-w-lg">
      <ModalHeader onClose={onClose}>Import Devices</ModalHeader>
      <ModalBody>
        <div className="space-y-6">
          {/* Step 1: Download Template */}
          <section>
            <h3 className="text-sm font-medium text-gray-900 mb-2">
              Step 1: Download Template
            </h3>
            <p className="text-sm text-gray-500 mb-3">
              Download our CSV template to ensure your data is formatted correctly.
            </p>
            <Button variant="outline" onClick={downloadImportTemplate} className="gap-2">
              <Download className="h-4 w-4" />
              Download CSV Template
            </Button>
          </section>

          {/* Step 2: Upload File */}
          <section>
            <h3 className="text-sm font-medium text-gray-900 mb-2">Step 2: Upload File</h3>
            <div
              className={cn(
                'border-2 border-dashed rounded-lg p-8 text-center transition-colors',
                dragOver
                  ? 'border-blue-500 bg-blue-50'
                  : file
                  ? 'border-green-300 bg-green-50'
                  : 'border-gray-300 hover:border-gray-400'
              )}
              onDragOver={(e) => {
                e.preventDefault();
                setDragOver(true);
              }}
              onDragLeave={() => setDragOver(false)}
              onDrop={handleDrop}
            >
              {file ? (
                <div className="flex items-center justify-center gap-3">
                  <File className="h-8 w-8 text-green-600" />
                  <div className="text-left">
                    <div className="font-medium text-gray-900">{file.name}</div>
                    <div className="text-sm text-gray-500">
                      {(file.size / 1024).toFixed(1)} KB
                    </div>
                  </div>
                  <button
                    type="button"
                    onClick={() => setFile(null)}
                    className="ml-4 text-gray-400 hover:text-gray-600"
                  >
                    <X className="h-5 w-5" />
                  </button>
                </div>
              ) : (
                <>
                  <Upload className="h-10 w-10 text-gray-400 mx-auto mb-3" />
                  <p className="text-gray-600 mb-2">Drag and drop your CSV file here</p>
                  <p className="text-sm text-gray-500 mb-3">or</p>
                  <Button variant="outline" onClick={() => fileInputRef.current?.click()}>
                    Browse Files
                  </Button>
                  <input
                    ref={fileInputRef}
                    type="file"
                    accept=".csv"
                    onChange={handleFileSelect}
                    className="hidden"
                  />
                </>
              )}
            </div>
          </section>

          {/* Step 3: Review & Import */}
          {importResult && (
            <section>
              <h3 className="text-sm font-medium text-gray-900 mb-2">Step 3: Results</h3>
              <div className="bg-gray-50 rounded-lg p-4 space-y-3">
                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-2 text-green-600">
                    <CheckCircle2 className="h-5 w-5" />
                    <span className="font-medium">{importResult.success} successful</span>
                  </div>
                  {importResult.failed > 0 && (
                    <div className="flex items-center gap-2 text-red-600">
                      <AlertCircle className="h-5 w-5" />
                      <span className="font-medium">{importResult.failed} failed</span>
                    </div>
                  )}
                </div>

                <div className="text-sm text-gray-600">
                  <span>{importResult.created} created</span>
                  {importResult.updated > 0 && (
                    <span> | {importResult.updated} updated</span>
                  )}
                </div>

                {importResult.errors.length > 0 && (
                  <>
                    <button
                      type="button"
                      onClick={() => setShowErrors(!showErrors)}
                      className="text-sm text-blue-600 hover:text-blue-800"
                    >
                      {showErrors ? 'Hide errors' : 'Show errors'}
                    </button>
                    {showErrors && (
                      <div className="mt-2 max-h-40 overflow-auto text-sm space-y-1">
                        {importResult.errors.map((error, i) => (
                          <div
                            key={i}
                            className="flex gap-2 text-red-600 bg-red-50 p-2 rounded"
                          >
                            <span className="font-mono">Row {error.row}:</span>
                            <span>{error.error}</span>
                          </div>
                        ))}
                      </div>
                    )}
                  </>
                )}
              </div>
            </section>
          )}
        </div>
      </ModalBody>
      <ModalFooter>
        <Button variant="outline" onClick={onClose} disabled={isLoading}>
          Cancel
        </Button>
        <Button onClick={handleImport} disabled={!file || isLoading || !!importResult}>
          {isLoading ? 'Importing...' : 'Import'}
        </Button>
      </ModalFooter>
    </Modal>
  );
}
